package rawrecording

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording/processing"
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
	CacheDir  string
	WorkDir   string

	// resolvedInputPath is the local path to the input (after S3 download if needed)
	resolvedInputPath string
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
	pf.String(FlagInputFile, "", DescInputFile)
	pf.String(FlagInputDir, "", DescInputDir)
	pf.String(FlagInputS3, "", DescInputS3)
	pf.String(FlagOutput, "", DescOutput)
	pf.Bool(FlagVerbose, false, DescVerbose)
	pf.String(FlagCacheDir, "", DescCacheDir)

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
	inputFile, _ := cmd.Flags().GetString(FlagInputFile)
	inputDir, _ := cmd.Flags().GetString(FlagInputDir)
	inputS3, _ := cmd.Flags().GetString(FlagInputS3)
	output, _ := cmd.Flags().GetString(FlagOutput)
	verbose, _ := cmd.Flags().GetBool(FlagVerbose)
	cacheDir, _ := cmd.Flags().GetString(FlagCacheDir)

	// Use default cache directory if not specified
	if cacheDir == "" {
		cacheDir = GetDefaultCacheDir()
	}

	return &GlobalArgs{
		InputFile: inputFile,
		InputDir:  inputDir,
		InputS3:   inputS3,
		Output:    output,
		Verbose:   verbose,
		CacheDir:  cacheDir,
	}, nil
}

// validateGlobalArgs validates global arguments
func validateGlobalArgs(globalArgs *GlobalArgs, requireOutput bool) error {
	if globalArgs.InputFile == "" && globalArgs.InputDir == "" && globalArgs.InputS3 == "" {
		return fmt.Errorf("either --%s or --%s or --%s must be specified", FlagInputFile, FlagInputDir, FlagInputS3)
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
		return fmt.Errorf("--%s, --%s and --%s are exclusive, only one is allowed", FlagInputFile, FlagInputDir, FlagInputS3)
	}

	if requireOutput && globalArgs.Output == "" {
		return fmt.Errorf("--%s directory must be specified", FlagOutput)
	}

	return nil
}

// resolveInputPath resolves the input to a local path, downloading from S3 if necessary
func resolveInputPath(ctx context.Context, globalArgs *GlobalArgs) (string, error) {
	// If already resolved, return cached path
	if globalArgs.resolvedInputPath != "" {
		return globalArgs.resolvedInputPath, nil
	}

	var inputPath string

	if globalArgs.InputFile != "" {
		inputPath = globalArgs.InputFile
	} else if globalArgs.InputDir != "" {
		inputPath = globalArgs.InputDir
	} else if globalArgs.InputS3 != "" {
		// Download from S3 (with caching)
		downloader := NewS3Downloader(globalArgs.CacheDir, globalArgs.Verbose)
		downloadedPath, err := downloader.Download(ctx, globalArgs.InputS3)
		if err != nil {
			return "", fmt.Errorf("failed to download from S3: %w", err)
		}
		inputPath = downloadedPath
	} else {
		return "", fmt.Errorf("no input specified")
	}

	// Cache the resolved path
	globalArgs.resolvedInputPath = inputPath
	return inputPath, nil
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
		return nil, fmt.Errorf("only one filter can be specified at a time: --%s, --%s, and --%s are mutually exclusive", FlagUserID, FlagSessionID, FlagTrackID)
	}

	// Resolve input path (download from S3 if needed)
	inputPath, err := resolveInputPath(context.Background(), globalArgs)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("%s '%s' not found in recording. Use 'list-tracks --%s tracks' to see available track IDs", FlagTrackID, trackID, FlagFormat)
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
			return nil, fmt.Errorf("%s '%s' not found in recording. Use 'list-tracks --%s sessions' to see available session IDs", FlagSessionID, sessionID, FlagFormat)
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
			return nil, fmt.Errorf("%s '%s' not found in recording. Use 'list-tracks --%s users' to see available user IDs", FlagUserID, userID, FlagFormat)
		}
	}

	return metadata, nil
}

// setupLogger creates a logger with the specified verbosity
func setupLogger(verbose bool) *processing.ProcessingLogger {
	var level getstream.LogLevel
	if verbose {
		level = getstream.LogLevelDebug
	} else {
		level = getstream.LogLevelInfo
	}
	return processing.NewRawToolLogger(getstream.NewDefaultLogger(os.Stderr, "", log.LstdFlags, level))
}

// prepareWorkDir extracts the recording to a temp directory and returns the working directory
func prepareWorkDir(globalArgs *GlobalArgs, logger *processing.ProcessingLogger) (string, func(), error) {
	// Resolve input path (download from S3 if needed)
	path, err := resolveInputPath(context.Background(), globalArgs)
	if err != nil {
		return "", nil, err
	}

	workingDir, cleanup, err := processing.ExtractToTempDir(path, logger)
	if err != nil {
		return "", nil, fmt.Errorf("failed to prepare working directory: %w", err)
	}

	return workingDir, cleanup, nil
}

// completeUserIDs provides completion for user IDs
func completeUserIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	inputFile, _ := cmd.Flags().GetString(FlagInputFile)
	inputDir, _ := cmd.Flags().GetString(FlagInputDir)

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
	inputFile, _ := cmd.Flags().GetString(FlagInputFile)
	inputDir, _ := cmd.Flags().GetString(FlagInputDir)

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
	inputFile, _ := cmd.Flags().GetString(FlagInputFile)
	inputDir, _ := cmd.Flags().GetString(FlagInputDir)

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
