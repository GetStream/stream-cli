package rawrecording

import (
	"fmt"
	"os"

	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func extractAudioCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract-audio",
		Short: "Generate playable audio files from raw recording tracks",
		Long: heredoc.Doc(`
			Generate playable audio files from raw recording tracks.

			Supports formats: webm, mp3, and others.

			Filters are mutually exclusive: you can only specify one of
			--user-id, --session-id, or --track-id at a time.
		`),
		Example: heredoc.Doc(`
			# Extract audio for all users (no filters)
			$ stream-cli video raw-recording extract-audio --input-file recording.zip --output ./out

			# Extract audio for specific user (all their tracks)
			$ stream-cli video raw-recording extract-audio --input-file recording.zip --output ./out --user-id user123

			# Extract audio for specific session
			$ stream-cli video raw-recording extract-audio --input-file recording.zip --output ./out --session-id session456

			# Extract a specific track
			$ stream-cli video raw-recording extract-audio --input-file recording.zip --output ./out --track-id track1
		`),
		RunE: runExtractAudio,
	}

	fl := cmd.Flags()
	fl.String("user-id", "", "Filter by user ID")
	fl.String("session-id", "", "Filter by session ID")
	fl.String("track-id", "", "Filter by track ID")
	fl.Bool("fill-gaps", true, "Fill with silence when track was muted")
	fl.Bool("fix-dtx", true, "Fix DTX shrink audio")

	// Register completions
	_ = cmd.RegisterFlagCompletionFunc("user-id", completeUserIDs)
	_ = cmd.RegisterFlagCompletionFunc("session-id", completeSessionIDs)
	_ = cmd.RegisterFlagCompletionFunc("track-id", completeTrackIDs)

	return cmd
}

func runExtractAudio(cmd *cobra.Command, args []string) error {
	globalArgs, err := getGlobalArgs(cmd)
	if err != nil {
		return err
	}

	// Validate global args (output is required for extract-audio)
	if err := validateGlobalArgs(globalArgs, true); err != nil {
		return err
	}

	userID, _ := cmd.Flags().GetString("user-id")
	sessionID, _ := cmd.Flags().GetString("session-id")
	trackID, _ := cmd.Flags().GetString("track-id")
	fillGaps, _ := cmd.Flags().GetBool("fill-gaps")
	fixDtx, _ := cmd.Flags().GetBool("fix-dtx")

	// Validate input arguments against actual recording data
	metadata, err := validateInputArgs(globalArgs, userID, sessionID, trackID)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	logger := setupLogger(globalArgs.Verbose)
	logger.Info("Starting extract-audio command")

	// Print banner
	printExtractAudioBanner(cmd, globalArgs, userID, sessionID, trackID, fillGaps, fixDtx)

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

	// Extract audio tracks
	if err := processing.ExtractTracks(globalArgs.WorkDir, globalArgs.Output, userID, sessionID, trackID, metadata, "audio", "both", fillGaps, fixDtx, logger); err != nil {
		return fmt.Errorf("failed to extract audio: %w", err)
	}

	logger.Info("Extract audio command completed")
	return nil
}

func printExtractAudioBanner(cmd *cobra.Command, globalArgs *GlobalArgs, userID, sessionID, trackID string, fillGaps, fixDtx bool) {
	cmd.Println("Extract audio command with mutually exclusive filtering:")
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
		cmd.Printf("  -> Processing all audio tracks for session '%s'\n", sessionID)
	} else if userID != "" {
		cmd.Printf("  -> Processing all audio tracks for user '%s'\n", userID)
	} else {
		cmd.Println("  -> Processing all audio tracks (no filters)")
	}
	cmd.Printf("  Fill gaps: %t\n", fillGaps)
	cmd.Printf("  Fix DTX: %t\n", fixDtx)
}
