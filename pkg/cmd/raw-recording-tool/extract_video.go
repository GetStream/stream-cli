package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GetStream/getstream-go/v3"
	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
)

type ExtractVideoArgs struct {
	UserID    string
	SessionID string
	TrackID   string
	FillGaps  bool
}

type ExtractVideoProcess struct {
	logger *getstream.DefaultLogger
}

func NewExtractVideoProcess(logger *getstream.DefaultLogger) *ExtractVideoProcess {
	return &ExtractVideoProcess{logger: logger}
}

func (p *ExtractVideoProcess) runExtractVideo(args []string, globalArgs *GlobalArgs) {
	printHelpIfAsked(args, p.printUsage)

	// Parse command-specific flags
	fs := flag.NewFlagSet("extract-video", flag.ExitOnError)
	extractVideoArgs := &ExtractVideoArgs{}
	fs.StringVar(&extractVideoArgs.UserID, "userId", "", "Specify a userId (empty for all)")
	fs.StringVar(&extractVideoArgs.SessionID, "sessionId", "", "Specify a sessionId (empty for all)")
	fs.StringVar(&extractVideoArgs.TrackID, "trackId", "", "Specify a trackId (empty for all)")
	fs.BoolVar(&extractVideoArgs.FillGaps, "fill_gaps", true, "Fill with black frame when track was muted (default true)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Validate input arguments against actual recording data
	metadata, err := validateInputArgs(globalArgs, extractVideoArgs.UserID, extractVideoArgs.SessionID, extractVideoArgs.TrackID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Validation error: %v\n", err)
		os.Exit(1)
	}

	p.logger.Info("Starting extract-video command")
	p.printBanner(globalArgs, extractVideoArgs)

	// Extract video tracks
	if e := extractVideoTracks(globalArgs, extractVideoArgs, metadata, p.logger); e != nil {
		p.logger.Error("Failed to extract video tracks: %v", e)
		os.Exit(1)
	}

	p.logger.Info("Extract video command completed successfully")
}

func (p *ExtractVideoProcess) printBanner(globalArgs *GlobalArgs, extractVideoArgs *ExtractVideoArgs) {
	fmt.Printf("Extract video command with mutually exclusive filtering:\n")
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
	fmt.Printf("  User ID filter: %s\n", extractVideoArgs.UserID)
	fmt.Printf("  Session ID filter: %s\n", extractVideoArgs.SessionID)
	fmt.Printf("  Track ID filter: %s\n", extractVideoArgs.TrackID)

	if extractVideoArgs.TrackID != "" {
		fmt.Printf("  → Processing specific track '%s'\n", extractVideoArgs.TrackID)
	} else if extractVideoArgs.SessionID != "" {
		fmt.Printf("  → Processing all video tracks for session '%s'\n", extractVideoArgs.SessionID)
	} else if extractVideoArgs.UserID != "" {
		fmt.Printf("  → Processing all video tracks for user '%s'\n", extractVideoArgs.UserID)
	} else {
		fmt.Printf("  → Processing all video tracks (no filters)\n")
	}
	fmt.Printf("  Fill gaps: %t\n", extractVideoArgs.FillGaps)
}

func (p *ExtractVideoProcess) printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: raw-tools [global options] extract-video [command options]\n\n")
	fmt.Fprintf(os.Stderr, "Generate playable video files from raw recording tracks.\n")
	fmt.Fprintf(os.Stderr, "Supports formats: webm, mp4, and others.\n\n")
	fmt.Fprintf(os.Stderr, "Command Options (Mutually Exclusive Filters):\n")
	fmt.Fprintf(os.Stderr, "  --userId <id>          Filter by user ID\n")
	fmt.Fprintf(os.Stderr, "  --sessionId <id>       Filter by session ID\n")
	fmt.Fprintf(os.Stderr, "  --trackId <id>         Filter by track ID\n")
	fmt.Fprintf(os.Stderr, "                         (specify at most one of the above)\n")
	fmt.Fprintf(os.Stderr, "  --fill_gaps            Fill with black frames when muted\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  # Extract video for all users (no filters)\n")
	fmt.Fprintf(os.Stderr, "  raw-tools --inputFile recording.zip --output ./out extract-video\n\n")
	fmt.Fprintf(os.Stderr, "  # Extract video for specific user (all their tracks)\n")
	fmt.Fprintf(os.Stderr, "  raw-tools --inputFile recording.zip --output ./out extract-video --userId user123\n\n")
	fmt.Fprintf(os.Stderr, "  # Extract video for specific session (all users in that session)\n")
	fmt.Fprintf(os.Stderr, "  raw-tools --inputFile recording.zip --output ./out extract-video --sessionId session456\n\n")
	fmt.Fprintf(os.Stderr, "  # Extract a specific track\n")
	fmt.Fprintf(os.Stderr, "  raw-tools --inputFile recording.zip --output ./out extract-video --trackId track1\n\n")
	fmt.Fprintf(os.Stderr, "Global Options: Use 'raw-tools --help' to see global options.\n")
}

func extractVideoTracks(globalArgs *GlobalArgs, extractVideoArgs *ExtractVideoArgs, metadata *processing.RecordingMetadata, logger *getstream.DefaultLogger) error {
	return processing.ExtractTracks(globalArgs.WorkDir, globalArgs.Output, extractVideoArgs.UserID, extractVideoArgs.SessionID, extractVideoArgs.TrackID, metadata, "video", "both", extractVideoArgs.FillGaps, false, logger)
}
