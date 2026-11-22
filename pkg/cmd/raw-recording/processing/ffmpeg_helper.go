package processing

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type FileOffset struct {
	Name   string
	Offset int64
}

// audioGapFillFilter re-times one mix input against its own container timestamps.
//
// A segment captured from RTP carries a timestamp hole wherever the sender used
// DTX or packets were lost. The container spans the hole correctly, but the Opus
// decoder emits a fixed number of samples per packet and nothing for a hole, so
// the decoded audio is shorter than the segment claims and everything after each
// hole plays early. The error accumulates: a 0.04% loss rate measured 3.3s of
// drift across one 2h23m segment, which desynchronised that speaker from the
// other by the same amount.
//
// async=1 makes swresample insert silence so the samples line up with the frame
// timestamps again. min_hard_comp is the smallest hole it will correct and
// defaults to 0.1s, which leaves the common single-packet loss untouched; 20ms is
// one Opus frame, the smallest hole that can occur.
const audioGapFillFilter = "aresample=async=1:min_hard_comp=0.020"

// mp3EncoderFrameFilter repacks the filtergraph output into whole MP3 frames
// (1152 samples) before it reaches libmp3lame.
//
// LAME's float entry point reads samples in blocks of 8, so libavcodec rejects
// any frame whose plane is not padded to FFALIGN(nb_samples, 8) —
// "inadequate AVFrame plane padding", which fails the whole encode with EINVAL.
// ffmpeg cuts the graph output into encoder-sized frames by advancing the plane
// pointer and shrinking linesize by the same amount, so a cut frame carries only
// the allocation slack of the frame it came from. apad emits its trailing
// silence as a single frame with as little as 2 samples of slack, and when the
// stream's last cut needs more than that, the final frame is rejected. Feeding
// the encoder frames it never has to cut removes the failure, and leaves the
// encoded bytes identical on the inputs that already worked.
//
// p=0 leaves the last frame short rather than zero-padding it to 1152 (the
// filter's default): padding feeds LAME real trailing silence that its gapless
// header then stops trimming, which would push the output past the recording
// window apad anchors it to.
const mp3EncoderFrameFilter = "asetnsamples=n=1152:p=0"

func generateConcatFileArguments(outputPath, concatPath string) []string {
	args := defaultArgs()
	args = append(args, "-f", "concat")
	args = append(args, "-safe", "0")
	args = append(args, "-i", concatPath)
	args = append(args, "-c", "copy")
	args = append(args, "-y", outputPath)
	return args
}

func generateMuxFilesArguments(fileName string, audioFile string, videoFile string, offsetMs float64) []string {
	args := defaultArgs()

	// Apply offset using itsoffset
	// If offset is positive (video ahead), delay audio
	// If offset is negative (audio ahead), delay video
	if offsetMs != 0 {
		offsetSeconds := offsetMs / 1000.0

		if offsetMs > 0 {
			// Video is ahead, delay audio
			args = append(args, "-itsoffset", fmt.Sprintf("%.3f", offsetSeconds))
			args = append(args, "-i", audioFile)
			args = append(args, "-i", videoFile)
		} else {
			args = append(args, "-i", audioFile)
			args = append(args, "-itsoffset", fmt.Sprintf("%.3f", -offsetSeconds))
			args = append(args, "-i", videoFile)
		}
	} else {
		args = append(args, "-i", audioFile)
		args = append(args, "-i", videoFile)
	}

	args = append(args, "-map", "0:a")
	args = append(args, "-map", "1:v")
	args = append(args, "-c", "copy")
	args = append(args, "-y", fileName)
	return args
}

// generateMixAudioFilesArguments builds the ffmpeg args that mix the given audio
// files. recordingDurationSeconds > 0 selects window mode: the (relative) offsets are re-based
// so the earliest input starts at leadingGapMs (the leading-silence gap), the mix
// is padded (apad) to recordingDurationSeconds for trailing silence, a single input is run
// through the filtergraph, and the result is re-encoded. With recordingDurationSeconds 0
// (legacy) the earliest input is normalized to zero and a single same-container
// file is stream-copied. apad is only emitted with recordingDurationSeconds > 0, never
// whole_dur=0, which pads infinitely on older ffmpeg (e.g. 4.4) and hangs.
func generateMixAudioFilesArguments(fileName, format string, files []*FileOffset, leadingGapMs int64, recordingDurationSeconds float64) []string {
	var filterParts []string
	var mixParts []string
	args := defaultArgs()
	args = append(args, "-filter_complex_threads", "1")

	sort.Slice(files, func(i, j int) bool {
		return files[i].Offset < files[j].Offset
	})

	windowed := recordingDurationSeconds > 0
	useFilter := len(files) > 1 || windowed

	var offsetToAdd int64
	for i, fo := range files {
		args = append(args, "-i", fo.Name)

		if useFilter {
			if i == 0 {
				offsetToAdd = leadingGapMs - fo.Offset
			}
			offset := fo.Offset + offsetToAdd

			// Close the input's internal timestamp holes before positioning it, so
			// adelay places audio that already spans its full recorded duration.
			chain := audioGapFillFilter
			if offset > 0 {
				// for stereo: offset|offset
				chain += fmt.Sprintf(",adelay=%d|%d", offset, offset)
			}
			label := fmt.Sprintf("a%d", i)
			filterParts = append(filterParts, fmt.Sprintf("[%d:a]%s[%s]", i, chain, label))
			mixParts = append(mixParts, fmt.Sprintf("[%s]", label))
		}
	}

	if useFilter {
		// Build amix filter
		filter := strings.Join(filterParts, "; ")
		if filter != "" {
			filter += "; "
		}
		filter += strings.Join(mixParts, "")
		if len(files) > 1 {
			filter += fmt.Sprintf("amix=inputs=%d:normalize=0", len(files))
		}
		// Pad the mix to the full recording window for trailing silence. A single
		// input has no amix, so apad also serves as its terminal filter.
		if windowed {
			if len(files) > 1 {
				filter += ","
			}
			filter += fmt.Sprintf("apad=whole_dur=%.3f", recordingDurationSeconds)
		}
		// Must stay terminal: any filter after it re-cuts the frames.
		if format == FormatMp3 {
			filter += "," + mp3EncoderFrameFilter
		}

		args = append(args, "-filter_complex", filter)
	}

	audioLib := audioLibForExtension(format)
	mkvAudioLib := audioLibForExtension(FormatMkv)
	// Copy is enough in case of webm, weba, mka, mkv when len == 1 and no window
	if windowed || audioLib != mkvAudioLib || len(files) > 1 {
		args = append(args, "-c:a", audioLib)
		if format == FormatMp3 {
			// Match the live GStreamer path (createMP3Output): lamemp3enc runs in
			// its default target=quality mode, i.e. VBR at quality 4 and 44100 Hz.
			// 44100 Hz is MP3-only: libopus supports only 8/12/16/24/48 kHz, so
			// pinning opus to 44100 would just force an implicit resample to 48k.
			args = append(args, "-q:a", "4")
			args = append(args, "-ar", "44100")
		} else {
			args = append(args, "-b:a", "128k")
		}
	} else {
		args = append(args, "-c", "copy")
	}

	if format == FormatWeba {
		args = append(args, "-f", "webm")
	}

	args = append(args, "-y", fileName)

	return args
}

func audioLibForExtension(str string) string {
	switch str {
	case FormatMp3:
		return "libmp3lame"
	case FormatWeba, FormatWebm, FormatMkv, FormatMka:
		return "libopus"
	default:
		return "libopus"
	}
}

// generateSilenceArguments builds ffmpeg args that write a silent audio file of
// the given duration in the container implied by format (the audio codec and,
// for weba, the output format are derived from it). It is used both for
// inter-segment gap fillers (format from the segment container extension) and
// for a full-window silence output when no participant audio was captured.
func generateSilenceArguments(fileName, format string, duration float64) []string {
	args := defaultArgs()
	args = append(args, "-f", "lavfi")
	args = append(args, "-t", fmt.Sprintf("%.3f", duration))
	args = append(args, "-i", "anullsrc=cl=stereo:r=48000")
	args = append(args, "-c:a", audioLibForExtension(format))
	args = append(args, "-b:a", "32k")
	if format == FormatWeba {
		args = append(args, "-f", "webm")
	}
	args = append(args, "-y", fileName)
	return args
}

func generateBlackVideoArguments(fileName, mimeType string, duration float64, width, height, frameRate int) []string {
	args := defaultArgs()
	lowerMimeType := strings.ToLower(mimeType)
	width = makeEven(width)
	height = makeEven(height)

	args = append(args, "-f", "lavfi")
	args = append(args, "-t", fmt.Sprintf("%.3f", duration))
	args = append(args, "-i", fmt.Sprintf("color=c=black:s=%dx%d:r=%d", width, height, frameRate))
	args = append(args, "-c:v", videoLibForMimeType(mimeType))

	if lowerMimeType == "video/h264" {
		args = append(args, "-preset", "ultrafast")
	} else {
		args = append(args, "-b:v", "0")
		args = append(args, "-cpu-used", "8")
	}

	args = append(args, "-crf", "45")
	args = append(args, "-y", fileName)
	return args
}

func makeEven(value int) int {
	if value%2 != 0 {
		return value + 1
	}
	return value
}

// generateTranscodeVideoArguments creates ffmpeg args to re-encode a video file to the target codec.
// Used when segments have different codecs due to a codec switch mid-recording.
func generateTranscodeVideoArguments(inputPath, outputPath, targetMimeType string) []string {
	args := defaultArgs()
	args = append(args, "-i", inputPath)
	args = append(args, "-c:v", videoLibForMimeType(targetMimeType))

	if strings.ToLower(targetMimeType) == "video/h264" {
		args = append(args, "-preset", "ultrafast")
	} else {
		args = append(args, "-b:v", "0")
		args = append(args, "-cpu-used", "4")
	}

	args = append(args, "-crf", "23")
	args = append(args, "-y", outputPath)
	return args
}

func videoLibForMimeType(str string) string {
	switch strings.ToLower(str) {
	case "video/vp8":
		return "libvpx"
	case "video/vp9":
		return "libvpx-vp9"
	case "video/h264":
		return "libx264"
	case "video/av1":
		return "libaom-av1"
	default:
		return "libvpx"
	}
}

func outputFormatForMimeType(str string) (extension, suffix string) {
	extension = mkvExtension
	suffix = mkvSuffix
	return extension, suffix
}

func defaultArgs() []string {
	var args []string
	args = append(args, "-hide_banner")
	args = append(args, "-threads", "1")
	args = append(args, "-filter_threads", "1")
	return args
}

func runFFmpegCommand(args []string, logger *ProcessingLogger) error {
	logger.Infof("Starting FFmpeg process with args: %s", args)

	startAt := time.Now()
	cmd := exec.Command("ffmpeg", args...)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	logger.Infof("FFmpeg process pid<%d> output:\n%s", pid(cmd), string(output))

	if err != nil {
		logger.Errorf("FFmpeg process pid<%d> failed: %v", pid(cmd), err)
		return fmt.Errorf("FFmpeg process pid<%d> failed in %s: %w", pid(cmd), time.Since(startAt).Round(time.Millisecond), err)
	}

	logger.Infof("FFmpeg process pid<%d> ended successfully in %s", pid(cmd), time.Since(startAt).Round(time.Millisecond))
	return nil
}

func pid(cmd *exec.Cmd) int {
	if cmd == nil || cmd.Process == nil {
		return 0
	}

	return cmd.Process.Pid
}
