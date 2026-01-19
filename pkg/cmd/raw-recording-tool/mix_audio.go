package main

import (
	"fmt"
	"os"

	"github.com/GetStream/getstream-go/v3"
	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
)

// MixAudioArgs represents the arguments for the mix-audio command
type MixAudioArgs struct {
	IncludeScreenShare bool
}

type MixAudioProcess struct {
	logger *getstream.DefaultLogger
}

func NewMixAudioProcess(logger *getstream.DefaultLogger) *MixAudioProcess {
	return &MixAudioProcess{logger: logger}
}

// runMixAudio handles the mix-audio command
func (p *MixAudioProcess) runMixAudio(args []string, globalArgs *GlobalArgs) {
	printHelpIfAsked(args, p.printUsage)

	mixAudioArgs := &MixAudioArgs{
		IncludeScreenShare: false,
	}

	// Validate input arguments against actual recording data
	metadata, err := validateInputArgs(globalArgs, "", "", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Validation error: %v\n", err)
		os.Exit(1)
	}

	p.logger.Info("Starting mix-audio command")

	// Execute the mix-audio operation
	if e := p.mixAllAudioTracks(globalArgs, mixAudioArgs, metadata, p.logger); e != nil {
		p.logger.Error("Mix-audio failed: %v", e)
		os.Exit(1)
	}

	p.logger.Info("Mix-audio command completed successfully")
}

// mixAllAudioTracks orchestrates the entire audio mixing workflow using existing extraction logic
func (p *MixAudioProcess) mixAllAudioTracks(globalArgs *GlobalArgs, mixAudioArgs *MixAudioArgs, metadata *processing.RecordingMetadata, logger *getstream.DefaultLogger) error {
	mixer := processing.NewAudioMixer(logger)
	mixer.MixAllAudioTracks(&processing.AudioMixerConfig{
		WorkDir:         globalArgs.WorkDir,
		OutputDir:       globalArgs.Output,
		WithScreenshare: false,
		WithExtract:     true,
		WithCleanup:     false,
	}, metadata, logger)
	return nil
}

// printMixAudioUsage prints the usage information for the mix-audio command
func (p *MixAudioProcess) printUsage() {
	fmt.Println("Usage: raw-tools [global-options] mix-audio [options]")
	fmt.Println()
	fmt.Println("Mix all audio tracks from multiple users/sessions into a single audio file")
	fmt.Println("with proper timing synchronization (like a conference call recording).")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --userId <id>        Filter by user ID (* for all users, default: *)")
	fmt.Println("  --sessionId <id>     Filter by session ID (* for all sessions, default: *)")
	fmt.Println("  --trackId <id>       Filter by track ID (* for all tracks, default: *)")
	fmt.Println("  --no-fill-gaps       Don't fill gaps with silence (not recommended for mixing)")
	fmt.Println("  -h, --help           Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Mix all audio tracks from all users and sessions")
	fmt.Println("  raw-tools --inputFile recording.tar.gz --output /tmp/mixed mix-audio")
	fmt.Println()
	fmt.Println("  # Mix audio tracks from a specific user")
	fmt.Println("  raw-tools --inputFile recording.tar.gz --output /tmp/mixed mix-audio --userId user123")
	fmt.Println()
	fmt.Println("  # Mix audio tracks from a specific session")
	fmt.Println("  raw-tools --inputFile recording.tar.gz --output /tmp/mixed mix-audio --sessionId session456")
	fmt.Println()
	fmt.Println("Output:")
	fmt.Println("  Creates 'mixed_audio.webm' - a single audio file containing all mixed tracks")
	fmt.Println("  with proper timing synchronization based on the original recording timeline.")
}
