package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GetStream/getstream-go/v3"
	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
)

type ProcessAllArgs struct {
	UserID    string
	SessionID string
	TrackID   string
}

type ProcessAllProcess struct {
	logger *getstream.DefaultLogger
}

func NewProcessAllProcess(logger *getstream.DefaultLogger) *ProcessAllProcess {
	return &ProcessAllProcess{logger: logger}
}

func (p *ProcessAllProcess) runProcessAll(args []string, globalArgs *GlobalArgs) {
	printHelpIfAsked(args, p.printUsage)

	// Parse command-specific flags
	fs := flag.NewFlagSet("process-all", flag.ExitOnError)
	processAllArgs := &ProcessAllArgs{}
	fs.StringVar(&processAllArgs.UserID, "userId", "", "Specify a userId (empty for all)")
	fs.StringVar(&processAllArgs.SessionID, "sessionId", "", "Specify a sessionId (empty for all)")
	fs.StringVar(&processAllArgs.TrackID, "trackId", "", "Specify a trackId (empty for all)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Validate input arguments against actual recording data
	metadata, err := validateInputArgs(globalArgs, processAllArgs.UserID, processAllArgs.SessionID, processAllArgs.TrackID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Validation error: %v\n", err)
		os.Exit(1)
	}

	p.logger.Info("Starting process-all command")

	// Display hierarchy information for user clarity
	fmt.Printf("Process-all command (audio + video + mux) with hierarchical filtering:\n")
	fmt.Printf("  Input file: %s\n", globalArgs.InputFile)
	fmt.Printf("  Output directory: %s\n", globalArgs.Output)
	fmt.Printf("  User ID filter: %s\n", processAllArgs.UserID)
	fmt.Printf("  Session ID filter: %s\n", processAllArgs.SessionID)
	fmt.Printf("  Track ID filter: %s\n", processAllArgs.TrackID)
	fmt.Printf("  Gap filling: always enabled\n")

	if processAllArgs.TrackID != "" {
		fmt.Printf("  → Processing specific track '%s'\n", processAllArgs.TrackID)
	} else if processAllArgs.SessionID != "" {
		fmt.Printf("  → Processing all tracks for session '%s'\n", processAllArgs.SessionID)
	} else if processAllArgs.UserID != "" {
		fmt.Printf("  → Processing all tracks for user '%s'\n", processAllArgs.UserID)
	} else {
		fmt.Printf("  → Processing all tracks (no filters)\n")
	}

	// Process all tracks and mux them
	if err := p.processAllTracks(globalArgs, processAllArgs, metadata, p.logger); err != nil {
		p.logger.Error("Failed to process and mux tracks: %v", err)
		os.Exit(1)
	}

	p.logger.Info("Process-all command completed successfully")
}

func (p *ProcessAllProcess) printUsage() {
	fmt.Printf("Usage: process-all [OPTIONS]\n")
	fmt.Printf("\nProcess audio, video, and mux them into combined files (all-in-one workflow)\n")
	fmt.Printf("Outputs 3 files per session: audio WebM, video WebM, and muxed WebM\n")
	fmt.Printf("Gap filling is always enabled for seamless playback.\n")
	fmt.Printf("\nOptions:\n")
	fmt.Printf("  --userId STRING    Specify a userId or * for all (default: \"*\")\n")
	fmt.Printf("  --sessionId STRING Specify a sessionId or * for all (default: \"*\")\n")
	fmt.Printf("  --trackId STRING   Specify a trackId or * for all (default: \"*\")\n")
	fmt.Printf("\nOutput files per session:\n")
	fmt.Printf("  audio_{userId}_{sessionId}_{trackId}.webm    - Audio-only file\n")
	fmt.Printf("  video_{userId}_{sessionId}_{trackId}.webm    - Video-only file\n")
	fmt.Printf("  muxed_{userId}_{sessionId}_{trackId}.webm    - Combined audio+video file\n")
}

func (p *ProcessAllProcess) processAllTracks(globalArgs *GlobalArgs, processAllArgs *ProcessAllArgs, metadata *processing.RecordingMetadata, logger *getstream.DefaultLogger) error {

	if e := processing.ExtractTracks(globalArgs.WorkDir, globalArgs.Output, "", "", "", metadata, "audio", "both", true, true, logger); e != nil {
		return e
	}

	if e := processing.ExtractTracks(globalArgs.WorkDir, globalArgs.Output, "", "", "", metadata, "video", "both", true, true, logger); e != nil {
		return e
	}

	mixer := processing.NewAudioMixer(logger)
	mixer.MixAllAudioTracks(&processing.AudioMixerConfig{
		WorkDir:         globalArgs.WorkDir,
		OutputDir:       globalArgs.Output,
		WithScreenshare: false,
		WithExtract:     false,
		WithCleanup:     false,
	}, metadata, logger)

	muxer := processing.NewAudioVideoMuxer(p.logger)
	if e := muxer.MuxAudioVideoTracks(&processing.AudioVideoMuxerConfig{
		WorkDir:     globalArgs.WorkDir,
		OutputDir:   globalArgs.Output,
		UserID:      "",
		SessionID:   "",
		TrackID:     "",
		Media:       "",
		WithExtract: false,
		WithCleanup: false,
	}, metadata, logger); e != nil {
		return e
	}

	return nil
}
