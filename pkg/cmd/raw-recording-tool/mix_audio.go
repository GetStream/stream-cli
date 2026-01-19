package rawrecording

import (
	"fmt"
	"os"

	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func mixAudioCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mix-audio",
		Short: "Mix multiple audio tracks into one file",
		Long: heredoc.Doc(`
			Mix all audio tracks from multiple users/sessions into a single audio file
			with proper timing synchronization (like a conference call recording).

			Creates 'mixed_audio.webm' - a single audio file containing all mixed tracks
			with proper timing synchronization based on the original recording timeline.
		`),
		Example: heredoc.Doc(`
			# Mix all audio tracks from all users and sessions
			$ stream-cli video raw-recording mix-audio --input-file recording.zip --output ./out

			# Mix with verbose logging
			$ stream-cli video raw-recording mix-audio --input-file recording.zip --output ./out --verbose
		`),
		RunE: runMixAudio,
	}

	return cmd
}

func runMixAudio(cmd *cobra.Command, args []string) error {
	globalArgs, err := getGlobalArgs(cmd)
	if err != nil {
		return err
	}

	// Validate global args (output is required for mix-audio)
	if err := validateGlobalArgs(globalArgs, true); err != nil {
		return err
	}

	// Validate input arguments against actual recording data
	metadata, err := validateInputArgs(globalArgs, "", "", "")
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	logger := setupLogger(globalArgs.Verbose)
	logger.Info("Starting mix-audio command")

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

	// Mix all audio tracks
	mixer := processing.NewAudioMixer(logger)
	mixer.MixAllAudioTracks(&processing.AudioMixerConfig{
		WorkDir:         globalArgs.WorkDir,
		OutputDir:       globalArgs.Output,
		WithScreenshare: false,
		WithExtract:     true,
		WithCleanup:     false,
	}, metadata, logger)

	logger.Info("Mix-audio command completed successfully")
	return nil
}
