package rawrecording

import (
	"fmt"
	"log"
	"os"

	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/getstream-go/v3"
)

// GlobalArgs holds the global arguments shared across all subcommands
type GlobalArgs struct {
	InputFile string
	InputDir  string
	InputS3   string
	Output    string
	Verbose   bool
	WorkDir   string
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raw-recording",
		Short: "Post-processing tools for raw video call recordings",
		Long: heredoc.Doc(`
			Post-processing tools for raw video call recordings.

			These commands allow you to extract, process, and mux audio/video
			tracks from raw recording archives.
		`),
		Example: heredoc.Doc(`
			# List all tracks in a recording
			$ stream-cli video raw-recording list-tracks --input-file recording.zip

			# Extract audio tracks for a specific user
			$ stream-cli video raw-recording extract-audio --input-file recording.zip --output ./out --user-id user123

			# Mux audio and video tracks
			$ stream-cli video raw-recording mux-av --input-file recording.zip --output ./out
		`),
	}

	// Persistent flags (global options available to all subcommands)
	pf := cmd.PersistentFlags()
	pf.String("input-file", "", "Raw recording zip file path")
	pf.String("input-dir", "", "Raw recording directory path")
	pf.String("input-s3", "", "Raw recording S3 path")
	pf.String("output", "", "Output directory")
	pf.Bool("verbose", false, "Enable verbose logging")

	// Add subcommands
	cmd.AddCommand(
		listTracksCmd(),
		extractAudioCmd(),
		extractVideoCmd(),
		muxAVCmd(),
		mixAudioCmd(),
		processAllCmd(),
	)

	return cmd
}

// getGlobalArgs extracts global arguments from cobra command flags
func getGlobalArgs(cmd *cobra.Command) (*GlobalArgs, error) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	inputDir, _ := cmd.Flags().GetString("input-dir")
	inputS3, _ := cmd.Flags().GetString("input-s3")
	output, _ := cmd.Flags().GetString("output")
	verbose, _ := cmd.Flags().GetBool("verbose")

	return &GlobalArgs{
		InputFile: inputFile,
		InputDir:  inputDir,
		InputS3:   inputS3,
		Output:    output,
		Verbose:   verbose,
	}, nil
}

// validateGlobalArgs validates global arguments
func validateGlobalArgs(globalArgs *GlobalArgs, requireOutput bool) error {
	if globalArgs.InputFile == "" && globalArgs.InputDir == "" && globalArgs.InputS3 == "" {
		return fmt.Errorf("either --input-file or --input-dir or --input-s3 must be specified")
	}

	num := 0
	if globalArgs.InputFile != "" {
		num++
	}
	if globalArgs.InputDir != "" {
		num++
	}
	if globalArgs.InputS3 != "" {
		num++
	}
	if num > 1 {
		return fmt.Errorf("--input-file, --input-dir and --input-s3 are exclusive, only one is allowed")
	}

	if requireOutput && globalArgs.Output == "" {
		return fmt.Errorf("--output directory must be specified")
	}

	return nil
}

// validateInputArgs validates input arguments using mutually exclusive logic
func validateInputArgs(globalArgs *GlobalArgs, userID, sessionID, trackID string) (*processing.RecordingMetadata, error) {
	// Count how many filters are specified
	filtersCount := 0
	if userID != "" {
		filtersCount++
	}
	if sessionID != "" {
		filtersCount++
	}
	if trackID != "" {
		filtersCount++
	}

	// Ensure filters are mutually exclusive
	if filtersCount > 1 {
		return nil, fmt.Errorf("only one filter can be specified at a time: --user-id, --session-id, and --track-id are mutually exclusive")
	}

	var inputPath string
	if globalArgs.InputFile != "" {
		inputPath = globalArgs.InputFile
	} else if globalArgs.InputDir != "" {
		inputPath = globalArgs.InputDir
	} else {
		return nil, fmt.Errorf("S3 input not implemented yet")
	}

	// Parse metadata to validate the single specified argument
	logger := setupLogger(false)
	parser := processing.NewMetadataParser(logger)
	metadata, err := parser.ParseMetadataOnly(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recording for validation: %w", err)
	}

	// If no filters specified, no validation needed
	if filtersCount == 0 {
		return metadata, nil
	}

	// Validate the single specified filter
	if trackID != "" {
		found := false
		for _, track := range metadata.Tracks {
			if track.TrackID == trackID {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("track-id '%s' not found in recording. Use 'list-tracks --format tracks' to see available track IDs", trackID)
		}
	} else if sessionID != "" {
		found := false
		for _, track := range metadata.Tracks {
			if track.SessionID == sessionID {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("session-id '%s' not found in recording. Use 'list-tracks --format sessions' to see available session IDs", sessionID)
		}
	} else if userID != "" {
		found := false
		for _, uid := range metadata.UserIDs {
			if uid == userID {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("user-id '%s' not found in recording. Use 'list-tracks --format users' to see available user IDs", userID)
		}
	}

	return metadata, nil
}

// setupLogger creates a logger with the specified verbosity
func setupLogger(verbose bool) *getstream.DefaultLogger {
	var level getstream.LogLevel
	if verbose {
		level = getstream.LogLevelDebug
	} else {
		level = getstream.LogLevelInfo
	}
	return getstream.NewDefaultLogger(os.Stderr, "", log.LstdFlags, level)
}

// prepareWorkDir extracts the recording to a temp directory and returns the working directory
func prepareWorkDir(globalArgs *GlobalArgs, logger *getstream.DefaultLogger) (string, func(), error) {
	path := globalArgs.InputFile
	if path == "" {
		path = globalArgs.InputDir
	}

	workingDir, cleanup, err := processing.ExtractToTempDir(path, logger)
	if err != nil {
		return "", nil, fmt.Errorf("failed to prepare working directory: %w", err)
	}

	return workingDir, cleanup, nil
}

// completeUserIDs provides completion for user IDs
func completeUserIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	inputDir, _ := cmd.Flags().GetString("input-dir")

	inputPath := inputFile
	if inputPath == "" {
		inputPath = inputDir
	}
	if inputPath == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	logger := setupLogger(false)
	parser := processing.NewMetadataParser(logger)
	metadata, err := parser.ParseMetadataOnly(inputPath)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return metadata.UserIDs, cobra.ShellCompDirectiveNoFileComp
}

// completeSessionIDs provides completion for session IDs
func completeSessionIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	inputDir, _ := cmd.Flags().GetString("input-dir")

	inputPath := inputFile
	if inputPath == "" {
		inputPath = inputDir
	}
	if inputPath == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	logger := setupLogger(false)
	parser := processing.NewMetadataParser(logger)
	metadata, err := parser.ParseMetadataOnly(inputPath)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return metadata.Sessions, cobra.ShellCompDirectiveNoFileComp
}

// completeTrackIDs provides completion for track IDs
func completeTrackIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	inputDir, _ := cmd.Flags().GetString("input-dir")

	inputPath := inputFile
	if inputPath == "" {
		inputPath = inputDir
	}
	if inputPath == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	logger := setupLogger(false)
	parser := processing.NewMetadataParser(logger)
	metadata, err := parser.ParseMetadataOnly(inputPath)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	trackIDs := make([]string, 0, len(metadata.Tracks))
	seen := make(map[string]bool)
	for _, track := range metadata.Tracks {
		if !seen[track.TrackID] {
			trackIDs = append(trackIDs, track.TrackID)
			seen[track.TrackID] = true
		}
	}

	return trackIDs, cobra.ShellCompDirectiveNoFileComp
}
