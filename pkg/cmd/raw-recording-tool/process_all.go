package rawrecording

import (
	"fmt"
	"os"

	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func processAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process-all",
		Short: "Process audio, video, and mux (all-in-one)",
		Long: heredoc.Doc(`
			Process audio, video, and mux them into combined files (all-in-one workflow).

			Outputs 3 files per session: audio WebM, video WebM, and muxed WebM.
			Gap filling is always enabled for seamless playback.

			Filters are mutually exclusive: you can only specify one of
			--user-id, --session-id, or --track-id at a time.

			Output files per session:
			  audio_{userId}_{sessionId}_{trackId}.webm    - Audio-only file
			  video_{userId}_{sessionId}_{trackId}.webm    - Video-only file
			  muxed_{userId}_{sessionId}_{trackId}.webm    - Combined audio+video file
		`),
		Example: heredoc.Doc(`
			# Process all tracks
			$ stream-cli video raw-recording process-all --input-file recording.zip --output ./out

			# Process tracks for specific user
			$ stream-cli video raw-recording process-all --input-file recording.zip --output ./out --user-id user123

			# Process tracks for specific session
			$ stream-cli video raw-recording process-all --input-file recording.zip --output ./out --session-id session456
		`),
		RunE: runProcessAll,
	}

	fl := cmd.Flags()
	fl.String("user-id", "", "Filter by user ID")
	fl.String("session-id", "", "Filter by session ID")
	fl.String("track-id", "", "Filter by track ID")

	// Register completions
	_ = cmd.RegisterFlagCompletionFunc("user-id", completeUserIDs)
	_ = cmd.RegisterFlagCompletionFunc("session-id", completeSessionIDs)
	_ = cmd.RegisterFlagCompletionFunc("track-id", completeTrackIDs)

	return cmd
}

func runProcessAll(cmd *cobra.Command, args []string) error {
	globalArgs, err := getGlobalArgs(cmd)
	if err != nil {
		return err
	}

	// Validate global args (output is required for process-all)
	if err := validateGlobalArgs(globalArgs, true); err != nil {
		return err
	}

	userID, _ := cmd.Flags().GetString("user-id")
	sessionID, _ := cmd.Flags().GetString("session-id")
	trackID, _ := cmd.Flags().GetString("track-id")

	// Validate input arguments against actual recording data
	metadata, err := validateInputArgs(globalArgs, userID, sessionID, trackID)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	logger := setupLogger(globalArgs.Verbose)
	logger.Info("Starting process-all command")

	// Print banner
	cmd.Println("Process-all command (audio + video + mux) with hierarchical filtering:")
	cmd.Printf("  Input file: %s\n", globalArgs.InputFile)
	cmd.Printf("  Output directory: %s\n", globalArgs.Output)
	cmd.Printf("  User ID filter: %s\n", userID)
	cmd.Printf("  Session ID filter: %s\n", sessionID)
	cmd.Printf("  Track ID filter: %s\n", trackID)
	cmd.Println("  Gap filling: always enabled")

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

	// Extract audio tracks
	if err := processing.ExtractTracks(globalArgs.WorkDir, globalArgs.Output, "", "", "", metadata, "audio", "both", true, true, logger); err != nil {
		return fmt.Errorf("failed to extract audio tracks: %w", err)
	}

	// Extract video tracks
	if err := processing.ExtractTracks(globalArgs.WorkDir, globalArgs.Output, "", "", "", metadata, "video", "both", true, true, logger); err != nil {
		return fmt.Errorf("failed to extract video tracks: %w", err)
	}

	// Mix all audio tracks
	mixer := processing.NewAudioMixer(logger)
	mixer.MixAllAudioTracks(&processing.AudioMixerConfig{
		WorkDir:         globalArgs.WorkDir,
		OutputDir:       globalArgs.Output,
		WithScreenshare: false,
		WithExtract:     false,
		WithCleanup:     false,
	}, metadata, logger)

	// Mux audio/video tracks
	muxer := processing.NewAudioVideoMuxer(logger)
	if err := muxer.MuxAudioVideoTracks(&processing.AudioVideoMuxerConfig{
		WorkDir:     globalArgs.WorkDir,
		OutputDir:   globalArgs.Output,
		UserID:      "",
		SessionID:   "",
		TrackID:     "",
		Media:       "",
		WithExtract: false,
		WithCleanup: false,
	}, metadata, logger); err != nil {
		return fmt.Errorf("failed to mux audio/video tracks: %w", err)
	}

	logger.Info("Process-all command completed successfully")
	return nil
}
