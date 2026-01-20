package rawrecording

import (
	"fmt"
	"os"

	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func extractVideoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract-video",
		Short: "Generate playable video files from raw recording tracks",
		Long: heredoc.Doc(`
			Generate playable video files from raw recording tracks.

			Supports formats: webm, mp4, and others.

			Filters are mutually exclusive: you can only specify one of
			--user-id, --session-id, or --track-id at a time.
		`),
		Example: heredoc.Doc(`
			# Extract video for all users (no filters)
			$ stream-cli video raw-recording extract-video --input-file recording.zip --output ./out

			# Extract video for specific user (all their tracks)
			$ stream-cli video raw-recording extract-video --input-file recording.zip --output ./out --user-id user123

			# Extract video for specific session
			$ stream-cli video raw-recording extract-video --input-file recording.zip --output ./out --session-id session456

			# Extract a specific track
			$ stream-cli video raw-recording extract-video --input-file recording.zip --output ./out --track-id track1
		`),
		RunE: runExtractVideo,
	}

	fl := cmd.Flags()
	fl.String(FlagUserID, "", DescUserID)
	fl.String(FlagSessionID, "", DescSessionID)
	fl.String(FlagTrackID, "", DescTrackID)
	fl.Bool(FlagFillGaps, true, DescFillGapsVideo)

	// Register completions
	_ = cmd.RegisterFlagCompletionFunc(FlagUserID, completeUserIDs)
	_ = cmd.RegisterFlagCompletionFunc(FlagSessionID, completeSessionIDs)
	_ = cmd.RegisterFlagCompletionFunc(FlagTrackID, completeTrackIDs)

	return cmd
}

func runExtractVideo(cmd *cobra.Command, args []string) error {
	globalArgs, err := getGlobalArgs(cmd)
	if err != nil {
		return err
	}

	// Validate global args (output is required for extract-video)
	if err := validateGlobalArgs(globalArgs, true); err != nil {
		return err
	}

	userID, _ := cmd.Flags().GetString(FlagUserID)
	sessionID, _ := cmd.Flags().GetString(FlagSessionID)
	trackID, _ := cmd.Flags().GetString(FlagTrackID)
	fillGaps, _ := cmd.Flags().GetBool(FlagFillGaps)

	// Validate input arguments against actual recording data
	metadata, err := validateInputArgs(globalArgs, userID, sessionID, trackID)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	logger := setupLogger(globalArgs.Verbose)
	logger.Info("Starting extract-video command")

	// Print banner
	printExtractVideoBanner(cmd, globalArgs, userID, sessionID, trackID, fillGaps)

	// Prepare working directory
	workDir, cleanup, err := prepareWorkDir(globalArgs, logger)
	if err != nil {
		return err
	}
	defer cleanup()
	globalArgs.WorkDir = workDir

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(globalArgs.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Extract video tracks
	if err := processing.ExtractTracks(globalArgs.WorkDir, globalArgs.Output, userID, sessionID, trackID, metadata, "video", "both", fillGaps, false, logger); err != nil {
		return fmt.Errorf("failed to extract video tracks: %w", err)
	}

	logger.Info("Extract video command completed successfully")
	return nil
}

func printExtractVideoBanner(cmd *cobra.Command, globalArgs *GlobalArgs, userID, sessionID, trackID string, fillGaps bool) {
	cmd.Println("Extract video command with mutually exclusive filtering:")
	if globalArgs.InputFile != "" {
		cmd.Printf("  Input file: %s\n", globalArgs.InputFile)
	}
	if globalArgs.InputDir != "" {
		cmd.Printf("  Input directory: %s\n", globalArgs.InputDir)
	}
	if globalArgs.InputS3 != "" {
		cmd.Printf("  Input S3: %s\n", globalArgs.InputS3)
	}
	cmd.Printf("  Output directory: %s\n", globalArgs.Output)
	cmd.Printf("  User ID filter: %s\n", userID)
	cmd.Printf("  Session ID filter: %s\n", sessionID)
	cmd.Printf("  Track ID filter: %s\n", trackID)

	if trackID != "" {
		cmd.Printf("  -> Processing specific track '%s'\n", trackID)
	} else if sessionID != "" {
		cmd.Printf("  -> Processing all video tracks for session '%s'\n", sessionID)
	} else if userID != "" {
		cmd.Printf("  -> Processing all video tracks for user '%s'\n", userID)
	} else {
		cmd.Println("  -> Processing all video tracks (no filters)")
	}
	cmd.Printf("  Fill gaps: %t\n", fillGaps)
}
