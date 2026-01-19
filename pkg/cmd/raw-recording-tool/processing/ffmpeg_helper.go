package processing

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/GetStream/getstream-go/v3"
)

const TmpDir = "/tmp"

type FileOffset struct {
	Name   string
	Offset int64
}

func concatFile(outputPath string, files []string, logger *getstream.DefaultLogger) error {
	// Write to a temporary file
	tmpFile, err := os.CreateTemp(TmpDir, "concat_*.txt")
	if err != nil {
		return err
	}
	defer func() {
		tmpFile.Close()
		//		_ = os.Remove(concatFile.Name())
	}()

	for _, file := range files {
		if _, err := tmpFile.WriteString(fmt.Sprintf("file '%s'\n", file)); err != nil {
			return err
		}
	}

	args := []string{}
	args = append(args, "-f", "concat")
	args = append(args, "-safe", "0")
	args = append(args, "-i", tmpFile.Name())
	args = append(args, "-c", "copy")
	args = append(args, outputPath)
	return runFFMEPGCpmmand(args, logger)
}

func muxFiles(fileName string, audioFile string, videoFile string, offsetMs float64, logger *getstream.DefaultLogger) error {
	args := []string{}

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
	args = append(args, fileName)

	return runFFMEPGCpmmand(args, logger)
}

func mixAudioFiles(fileName string, files []*FileOffset, logger *getstream.DefaultLogger) error {
	var args []string

	var filterParts []string
	var mixParts []string

	sort.Slice(files, func(i, j int) bool {
		return files[i].Offset < files[j].Offset
	})

	var offsetToAdd int64
	for i, fo := range files {
		args = append(args, "-i", fo.Name)

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

	// Build amix filter
	filter := strings.Join(filterParts, "; ")
	if filter != "" {
		filter += "; "
	}
	filter += strings.Join(mixParts, "") +
		fmt.Sprintf("amix=inputs=%d:normalize=0", len(files))

	args = append(args, "-filter_complex", filter)
	args = append(args, "-c:a", "libopus")
	args = append(args, "-b:a", "128k")
	args = append(args, fileName)

	fmt.Println(strings.Join(args, " "))

	return runFFMEPGCpmmand(args, logger)
}

func generateSilence(fileName string, duration float64, logger *getstream.DefaultLogger) error {
	args := []string{}
	args = append(args, "-f", "lavfi")
	args = append(args, "-t", fmt.Sprintf("%.3f", duration))
	args = append(args, "-i", "anullsrc=cl=stereo:r=48000")
	args = append(args, "-c:a", "libopus")
	args = append(args, "-b:a", "32k")
	args = append(args, fileName)

	return runFFMEPGCpmmand(args, logger)
}

func generateBlackVideo(fileName, mimeType string, duration float64, width, height, frameRate int, logger *getstream.DefaultLogger) error {
	var codecLib string
	switch strings.ToLower(mimeType) {
	case "video/vp8":
		codecLib = "libvpx-vp9"
	case "video/vp9":
		codecLib = "libvpx-vp9"
	case "video/h264":
		codecLib = "libh264"
	case "video/av1":
		codecLib = "libav1"
	}

	args := []string{}
	args = append(args, "-f", "lavfi")
	args = append(args, "-t", fmt.Sprintf("%.3f", duration))
	args = append(args, "-i", fmt.Sprintf("color=c=black:s=%dx%d:r=%d", width, height, frameRate))
	args = append(args, "-c:v", codecLib)
	args = append(args, "-b:v", "1M")
	args = append(args, fileName)

	return runFFMEPGCpmmand(args, logger)
}

func runFFMEPGCpmmand(args []string, logger *getstream.DefaultLogger) error {
	cmd := exec.Command("ffmpeg", args...)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("FFmpeg command failed: %v", err)
		logger.Error("FFmpeg output: %s", string(output))
		return fmt.Errorf("ffmpeg command failed: %w", err)
	}

	logger.Info("Successfully ran ffmpeg: %s", args)
	logger.Debug("FFmpeg output: %s", string(output))
	return nil
}
