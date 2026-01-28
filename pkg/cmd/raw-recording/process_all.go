package rawrecording

import (
	"fmt"
	"os"

	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording/processing"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func processAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process-all",
		Short: "Process audio, video, and mux (all-in-one)",
		Long: heredoc.Doc(`
			Process audio, video, and mux them into combined files (all-in-one workflow).

			Outputs multiple MKV files: individual audio/video tracks, muxed A/V, and mixed audio.
			Gap filling and DTX fix are always enabled for seamless playback and proper A/V sync.

			Filters are mutually exclusive: you can only specify one of
			--user-id, --session-id, or --track-id at a time.

			Output files:
			  individual_{callType}_{callId}_{userId}_{sessionId}_audio_only_{timestamp}.mkv
			  individual_{callType}_{callId}_{userId}_{sessionId}_video_only_{timestamp}.mkv
			  individual_{callType}_{callId}_{userId}_{sessionId}_audio_video_{timestamp}.mkv
			  composite_{callType}_{callId}_audio_{timestamp}.mkv
		`),
		Example: heredoc.Doc(`
			# Process all tracks
			$ stream-cli video raw-recording process-all --input-file recording.tar.gz --output ./out

			# Process tracks for specific user
			$ stream-cli video raw-recording process-all --input-file recording.tar.gz --output ./out --user-id user123

			# Process tracks for specific session
			$ stream-cli video raw-recording process-all --input-file recording.tar.gz --output ./out --session-id session456
		`),
		RunE: runProcessAll,
	}

	fl := cmd.Flags()
	fl.String(FlagUserID, "", DescUserID)
	fl.String(FlagSessionID, "", DescSessionID)
	fl.String(FlagTrackID, "", DescTrackID)

	// Register completions
	_ = cmd.RegisterFlagCompletionFunc(FlagUserID, completeUserIDs)
	_ = cmd.RegisterFlagCompletionFunc(FlagSessionID, completeSessionIDs)
	_ = cmd.RegisterFlagCompletionFunc(FlagTrackID, completeTrackIDs)

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

	userID, _ := cmd.Flags().GetString(FlagUserID)
	sessionID, _ := cmd.Flags().GetString(FlagSessionID)
	trackID, _ := cmd.Flags().GetString(FlagTrackID)

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
	extractor := processing.NewTrackExtractor(logger)
	if _, err := extractor.ExtractTracks(&processing.TrackExtractorConfig{
		WorkDir:   globalArgs.WorkDir,
		OutputDir: globalArgs.Output,
		UserID:    "",
		SessionID: "",
		TrackID:   "",
		TrackKind: TrackTypeAudio,
		MediaType: "both",
		FillGap:   true,
		FillDtx:   true,
	}, metadata); err != nil {
		return fmt.Errorf("failed to extract audio tracks: %w", err)
	}

	// Extract video tracks
	if _, err := extractor.ExtractTracks(&processing.TrackExtractorConfig{
		WorkDir:   globalArgs.WorkDir,
		OutputDir: globalArgs.Output,
		UserID:    "",
		SessionID: "",
		TrackID:   "",
		TrackKind: TrackTypeVideo,
		MediaType: "both",
		FillGap:   true,
		FillDtx:   true,
	}, metadata); err != nil {
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
	}, metadata)

	// Mux audio/video tracks
	muxer := processing.NewAudioVideoMuxer(logger)
	if _, err := muxer.MuxAudioVideoTracks(&processing.AudioVideoMuxerConfig{
		WorkDir:     globalArgs.WorkDir,
		OutputDir:   globalArgs.Output,
		UserID:      "",
		SessionID:   "",
		TrackID:     "",
		MediaType:   "",
		WithExtract: false,
		WithCleanup: false,
	}, metadata); err != nil {
		return fmt.Errorf("failed to mux audio/video tracks: %w", err)
	}

	logger.Info("Process-all command completed successfully")
	return nil
}
