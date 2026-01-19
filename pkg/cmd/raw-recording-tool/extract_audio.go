package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GetStream/getstream-go/v3"
	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
)

type ExtractAudioArgs struct {
	UserID    string
	SessionID string
	TrackID   string
	FillGaps  bool
	FixDtx    bool
}

type ExtractAudioProcess struct {
	logger *getstream.DefaultLogger
}

func NewExtractAudioProcess(logger *getstream.DefaultLogger) *ExtractAudioProcess {
	return &ExtractAudioProcess{logger: logger}
}

func (p *ExtractAudioProcess) runExtractAudio(args []string, globalArgs *GlobalArgs) {
	printHelpIfAsked(args, p.printUsage)

	// Parse command-specific flags
	fs := flag.NewFlagSet("extract-audio", flag.ExitOnError)
	extractAudioArgs := &ExtractAudioArgs{}
	fs.StringVar(&extractAudioArgs.UserID, "userId", "", "Specify a userId (empty for all)")
	fs.StringVar(&extractAudioArgs.SessionID, "sessionId", "", "Specify a sessionId (empty for all)")
	fs.StringVar(&extractAudioArgs.TrackID, "trackId", "", "Specify a trackId (empty for all)")
	fs.BoolVar(&extractAudioArgs.FillGaps, "fill_gaps", true, "Fill with silence when track was muted (default true)")
	fs.BoolVar(&extractAudioArgs.FixDtx, "fix_dtx", true, "Fix DTX shrink audio (default true)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Validate input arguments against actual recording data
	metadata, err := validateInputArgs(globalArgs, extractAudioArgs.UserID, extractAudioArgs.SessionID, extractAudioArgs.TrackID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Validation error: %v\n", err)
		os.Exit(1)
	}

	p.logger.Info("Starting extract-audio command")
	p.printBanner(globalArgs, extractAudioArgs)

	// Implement extract audio functionality
	if e := extractAudioTracks(globalArgs, extractAudioArgs, metadata, p.logger); e != nil {
		p.logger.Error("Failed to extract audio: %v", e)
	}

	p.logger.Info("Extract audio command completed")
}

func (p *ExtractAudioProcess) printBanner(globalArgs *GlobalArgs, extractAudioArgs *ExtractAudioArgs) {
	fmt.Printf("Extract audio command with mutually exclusive filtering:\n")
	if globalArgs.InputFile != "" {
		fmt.Printf("  Input file: %s\n", globalArgs.InputFile)
	}
	if globalArgs.InputDir != "" {
		fmt.Printf("  Input directory: %s\n", globalArgs.InputDir)
	}
	if globalArgs.InputS3 != "" {
		fmt.Printf("  Input S3: %s\n", globalArgs.InputS3)
	}
	fmt.Printf("  Output directory: %s\n", globalArgs.Output)
	fmt.Printf("  User ID filter: %s\n", extractAudioArgs.UserID)
	fmt.Printf("  Session ID filter: %s\n", extractAudioArgs.SessionID)
	fmt.Printf("  Track ID filter: %s\n", extractAudioArgs.TrackID)

	if extractAudioArgs.TrackID != "" {
		fmt.Printf("  → Processing specific track '%s'\n", extractAudioArgs.TrackID)
	} else if extractAudioArgs.SessionID != "" {
		fmt.Printf("  → Processing all audio tracks for session '%s'\n", extractAudioArgs.SessionID)
	} else if extractAudioArgs.UserID != "" {
		fmt.Printf("  → Processing all audio tracks for user '%s'\n", extractAudioArgs.UserID)
	} else {
		fmt.Printf("  → Processing all audio tracks (no filters)\n")
	}
	fmt.Printf("  Fill gaps: %t\n", extractAudioArgs.FillGaps)
	fmt.Printf("  Fix DTX: %t\n", extractAudioArgs.FixDtx)
}

func (p *ExtractAudioProcess) printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: raw-tools [global options] extract-audio [command options]\n\n")
	fmt.Fprintf(os.Stderr, "Generate playable audio files from raw recording tracks.\n")
	fmt.Fprintf(os.Stderr, "Supports formats: webm, mp3, and others.\n\n")
	fmt.Fprintf(os.Stderr, "Command Options (Mutually Exclusive Filters):\n")
	fmt.Fprintf(os.Stderr, "  --userId <id>          Filter by user ID\n")
	fmt.Fprintf(os.Stderr, "  --sessionId <id>       Filter by session ID\n")
	fmt.Fprintf(os.Stderr, "  --trackId <id>         Filter by track ID\n")
	fmt.Fprintf(os.Stderr, "                         (specify at most one of the above)\n")
	fmt.Fprintf(os.Stderr, "  --fill_gaps            Fix DTX shrink audio, fill with silence when muted\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  # Extract audio for all users (no filters)\n")
	fmt.Fprintf(os.Stderr, "  raw-tools --inputFile recording.zip --output ./out extract-audio\n\n")
	fmt.Fprintf(os.Stderr, "  # Extract audio for specific user (all their tracks)\n")
	fmt.Fprintf(os.Stderr, "  raw-tools --inputFile recording.zip --output ./out extract-audio --userId user123\n\n")
	fmt.Fprintf(os.Stderr, "  # Extract audio for specific session (all users in that session)\n")
	fmt.Fprintf(os.Stderr, "  raw-tools --inputFile recording.zip --output ./out extract-audio --sessionId session456\n\n")
	fmt.Fprintf(os.Stderr, "  # Extract a specific track\n")
	fmt.Fprintf(os.Stderr, "  raw-tools --inputFile recording.zip --output ./out extract-audio --trackId track1\n\n")
	fmt.Fprintf(os.Stderr, "Global Options: Use 'raw-tools --help' to see global options.\n")
}

func extractAudioTracks(globalArgs *GlobalArgs, extractAudioArgs *ExtractAudioArgs, metadata *processing.RecordingMetadata, logger *getstream.DefaultLogger) error {
	return processing.ExtractTracks(globalArgs.WorkDir, globalArgs.Output, extractAudioArgs.UserID, extractAudioArgs.SessionID, extractAudioArgs.TrackID, metadata, "audio", "both", extractAudioArgs.FillGaps, extractAudioArgs.FixDtx, logger)
}
