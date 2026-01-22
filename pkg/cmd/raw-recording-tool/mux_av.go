package rawrecording

import (
	"fmt"
	"os"

	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func muxAVCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mux-av",
		Short: "Mux audio and video tracks into a single file",
		Long: heredoc.Doc(`
			Mux audio and video tracks into a single file.

			This command combines audio and video tracks from the same
			user/session into a single playable file.

			Filters are mutually exclusive: you can only specify one of
			--user-id, --session-id, or --track-id at a time.

			Media filtering:
			  --media user     Only mux user camera audio/video pairs
			  --media display  Only mux display sharing audio/video pairs
			  --media both     Mux both types (default)
		`),
		Example: heredoc.Doc(`
			# Mux all tracks
			$ stream-cli video raw-recording mux-av --input-file recording.zip --output ./out

			# Mux tracks for specific user
			$ stream-cli video raw-recording mux-av --input-file recording.zip --output ./out --user-id user123

			# Mux only user camera tracks
			$ stream-cli video raw-recording mux-av --input-file recording.zip --output ./out --media user

			# Mux only display sharing tracks
			$ stream-cli video raw-recording mux-av --input-file recording.zip --output ./out --media display
		`),
		RunE: runMuxAV,
	}

	fl := cmd.Flags()
	fl.String(FlagUserID, "", DescUserID)
	fl.String(FlagSessionID, "", DescSessionID)
	fl.String(FlagTrackID, "", DescTrackID)
	fl.String(FlagMedia, DefaultMedia, DescMedia)

	// Register completions
	_ = cmd.RegisterFlagCompletionFunc(FlagUserID, completeUserIDs)
	_ = cmd.RegisterFlagCompletionFunc(FlagSessionID, completeSessionIDs)
	_ = cmd.RegisterFlagCompletionFunc(FlagTrackID, completeTrackIDs)
	_ = cmd.RegisterFlagCompletionFunc(FlagMedia, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{MediaUser, MediaDisplay, MediaBoth}, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func runMuxAV(cmd *cobra.Command, args []string) error {
	globalArgs, err := getGlobalArgs(cmd)
	if err != nil {
		return err
	}

	// Validate global args (output is required for mux-av)
	if err := validateGlobalArgs(globalArgs, true); err != nil {
		return err
	}

	userID, _ := cmd.Flags().GetString(FlagUserID)
	sessionID, _ := cmd.Flags().GetString(FlagSessionID)
	trackID, _ := cmd.Flags().GetString(FlagTrackID)
	media, _ := cmd.Flags().GetString(FlagMedia)

	// Validate input arguments against actual recording data
	metadata, err := validateInputArgs(globalArgs, userID, sessionID, trackID)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	logger := setupLogger(globalArgs.Verbose)
	logger.Info("Starting mux-av command")

	// Print banner
	cmd.Println("Mux audio and video command with hierarchical filtering:")
	cmd.Printf("  Input file: %s\n", globalArgs.InputFile)
	cmd.Printf("  Output directory: %s\n", globalArgs.Output)
	cmd.Printf("  User ID filter: %s\n", userID)
	cmd.Printf("  Session ID filter: %s\n", sessionID)
	cmd.Printf("  Track ID filter: %s\n", trackID)
	cmd.Printf("  Media filter: %s\n", media)

	if trackID != "" {
		cmd.Printf("  -> Processing specific track '%s'\n", trackID)
	} else if sessionID != "" {
		cmd.Printf("  -> Processing all tracks for session '%s'\n", sessionID)
	} else if userID != "" {
		cmd.Printf("  -> Processing all tracks for user '%s'\n", userID)
	} else {
		cmd.Println("  -> Processing all tracks (no filters)")
	}

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

	// Mux audio/video tracks
	muxer := processing.NewAudioVideoMuxer(logger)
	if _, err := muxer.MuxAudioVideoTracks(&processing.AudioVideoMuxerConfig{
		WorkDir:     globalArgs.WorkDir,
		OutputDir:   globalArgs.Output,
		UserID:      userID,
		SessionID:   sessionID,
		TrackID:     trackID,
		MediaType:   media,
		WithExtract: true,
		WithCleanup: false,
	}, metadata); err != nil {
		return fmt.Errorf("failed to mux audio/video tracks: %w", err)
	}

	logger.Info("Mux audio and video command completed successfully")
	return nil
}
