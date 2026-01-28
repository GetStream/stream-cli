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

func generateConcatFileArguments(outputPath, concatPath string) ([]string, error) {
	args := defaultArgs()
	args = append(args, "-f", "concat")
	args = append(args, "-safe", "0")
	args = append(args, "-i", concatPath)
	args = append(args, "-c", "copy")
	args = append(args, "-y", outputPath)
	return args, nil
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

func generateMixAudioFilesArguments(fileName, format string, files []*FileOffset) []string {
	var filterParts []string
	var mixParts []string
	args := defaultArgs()

	sort.Slice(files, func(i, j int) bool {
		return files[i].Offset < files[j].Offset
	})

	var offsetToAdd int64
	for i, fo := range files {
		args = append(args, "-i", fo.Name)

		if len(files) > 1 {
			if i == 0 {
				offsetToAdd = -fo.Offset
			}
			offset := fo.Offset + offsetToAdd

			if offset > 0 {
				// for stereo: offset|offset
				label := fmt.Sprintf("a%d", i)
				filterParts = append(filterParts,
					fmt.Sprintf("[%d:a]adelay=%d|%d[%s]", i, offset, offset, label))
				mixParts = append(mixParts, fmt.Sprintf("[%s]", label))
			} else {
				mixParts = append(mixParts, fmt.Sprintf("[%d:a]", i))
			}
		}
	}

	if len(files) > 1 {
		// Build amix filter
		filter := strings.Join(filterParts, "; ")
		if filter != "" {
			filter += "; "
		}
		filter += strings.Join(mixParts, "") +
			fmt.Sprintf("amix=inputs=%d:normalize=0", len(files))

		args = append(args, "-filter_complex", filter)
	}

	audioLib := audioLibForExtension(format)
	mkvAudioLib := audioLibForExtension(FormatMkv)
	// Copy is enough in case of webm, weba, mka, mkv when len == 1
	if audioLib != mkvAudioLib || len(files) > 1 {
		args = append(args, "-c:a", audioLibForExtension(format))
		args = append(args, "-b:a", "128k")
	} else {
		args = append(args, "-c", "copy")
	}

	if format == FormatWeba {
		args = append(args, "-f", "webm")
	}

	args = append(args, "-y", fileName)

	fmt.Println(strings.Join(args, " "))
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

func generateSilenceArguments(fileName string, duration float64) []string {
	args := defaultArgs()
	args = append(args, "-f", "lavfi")
	args = append(args, "-t", fmt.Sprintf("%.3f", duration))
	args = append(args, "-i", "anullsrc=cl=stereo:r=48000")
	args = append(args, "-c:a", "libopus")
	args = append(args, "-b:a", "32k")
	args = append(args, "-y", fileName)
	return args
}

func generateBlackVideoArguments(fileName, mimeType string, duration float64, width, height, frameRate int) []string {
	args := defaultArgs()
	args = append(args, "-f", "lavfi")
	args = append(args, "-t", fmt.Sprintf("%.3f", duration))
	args = append(args, "-i", fmt.Sprintf("color=c=black:s=%dx%d:r=%d", width, height, frameRate))
	args = append(args, "-c:v", videoLibForMimeType(mimeType))

	if strings.ToLower(mimeType) == "video/h264" {
		args = append(args, "-preset", "ultrafast")
	} else {
		args = append(args, "-b:v", "0")
		args = append(args, "-cpu-used", "8")
	}

	args = append(args, "-crf", "45")
	args = append(args, "-y", fileName)
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
	return
}

func defaultArgs() []string {
	var args []string
	args = append(args, "-hide_banner")
	args = append(args, "-threads", "1")
	args = append(args, "-filter_threads", "1")
	return args
}

func runFFmpegCommand(args []string, logger *ProcessingLogger) error {
	startAt := time.Now()
	cmd := exec.Command("ffmpeg", args...)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	logger.Infof("FFmpeg process pid<%d> with args: %s", cmd.Process.Pid, args)
	logger.Infof("FFmpeg process pid<%d> output:\n%s", cmd.Process.Pid, string(output))

	if err != nil {
		logger.Errorf("FFmpeg process pid<%d> failed: %v", cmd.Process.Pid, err)
		return fmt.Errorf("FFmpeg process pid<%d> failed in %s: %w", cmd.Process.Pid, time.Now().Sub(startAt).Round(time.Millisecond), err)
	}

	logger.Infof("FFmpeg process pid<%d> ended successfully in %s", cmd.Process.Pid, time.Now().Sub(startAt).Round(time.Millisecond))
	return nil
}
