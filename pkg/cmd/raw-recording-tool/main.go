package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GetStream/getstream-go/v3"
	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
)

type GlobalArgs struct {
	InputFile string
	InputDir  string
	InputS3   string
	Output    string
	Verbose   bool

	WorkDir string
}

func main() {
	if len(os.Args) < 2 {
		printGlobalUsage()
		os.Exit(1)
	}

	// Parse global flags first
	globalArgs := &GlobalArgs{}
	command, remainingArgs := parseGlobalFlags(os.Args[1:], globalArgs)

	if command == "" {
		printGlobalUsage()
		os.Exit(1)
	}

	// Setup logger
	logger := setupLogger(globalArgs.Verbose)

	switch command {
	case "list-tracks":
		p := NewListTracksProcess(logger)
		p.runListTracks(remainingArgs, globalArgs)
	case "completion":
		runCompletion(remainingArgs)
	case "help", "-h", "--help":
		printGlobalUsage()
	default:
		if e := processCommand(command, globalArgs, remainingArgs, logger); e != nil {
			logger.Error("Error processing command %s - %v", command, e)
			os.Exit(1)
		}
	}
}

func processCommand(command string, globalArgs *GlobalArgs, remainingArgs []string, logger *getstream.DefaultLogger) error {
	// Extract to temp directory if needed (unified approach)
	path := globalArgs.InputFile
	if path == "" {
		path = globalArgs.InputDir
	}

	workingDir, cleanup, err := processing.ExtractToTempDir(path, logger)
	if err != nil {
		return fmt.Errorf("failed to prepare working directory: %w", err)
	}
	defer cleanup()
	globalArgs.WorkDir = workingDir

	// Create output directory if it doesn't exist
	if e := os.MkdirAll(globalArgs.Output, 0755); e != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	switch command {
	case "extract-audio":
		p := NewExtractAudioProcess(logger)
		p.runExtractAudio(remainingArgs, globalArgs)
	case "extract-video":
		p := NewExtractVideoProcess(logger)
		p.runExtractVideo(remainingArgs, globalArgs)
	case "mux-av":
		p := NewMuxAudioVideoProcess(logger)
		p.runMuxAV(remainingArgs, globalArgs)
	case "mix-audio":
		p := NewMixAudioProcess(logger)
		p.runMixAudio(remainingArgs, globalArgs)
	case "process-all":
		p := NewProcessAllProcess(logger)
		p.runProcessAll(remainingArgs, globalArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printGlobalUsage()
		os.Exit(1)
	}

	return nil
}

// parseGlobalFlags parses global flags and returns the command and remaining args
func parseGlobalFlags(args []string, globalArgs *GlobalArgs) (string, []string) {
	fs := flag.NewFlagSet("global", flag.ContinueOnError)

	fs.StringVar(&globalArgs.InputFile, "inputFile", "", "Specify raw recording zip file on file system")
	fs.StringVar(&globalArgs.InputDir, "inputDir", "", "Specify raw recording directory on file system")
	fs.StringVar(&globalArgs.InputS3, "inputS3", "", "Specify raw recording zip file on S3")
	fs.StringVar(&globalArgs.Output, "output", "", "Specify an output directory")
	fs.BoolVar(&globalArgs.Verbose, "verbose", false, "Enable verbose logging")

	// Find the command by looking for known commands
	knownCommands := map[string]bool{
		"list-tracks":   true,
		"extract-audio": true,
		"extract-video": true,
		"mux-av":        true,
		"mix-audio":     true,
		"process-all":   true,
		"completion":    true,
		"help":          true,
	}

	commandIndex := -1
	for i, arg := range args {
		if knownCommands[arg] {
			commandIndex = i
			break
		}
	}

	if commandIndex == -1 {
		return "", nil
	}

	// Parse global flags (everything before the command)
	globalFlags := args[:commandIndex]
	command := args[commandIndex]
	remainingArgs := args[commandIndex+1:]

	err := fs.Parse(globalFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing global flags: %v\n", err)
		os.Exit(1)
	}

	// Validate global arguments
	if e := validateGlobalArgs(globalArgs, command); e != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		printGlobalUsage()
		os.Exit(1)
	}

	return command, remainingArgs
}

func setupLogger(verbose bool) *getstream.DefaultLogger {
	var level getstream.LogLevel
	if verbose {
		level = getstream.LogLevelDebug
	} else {
		level = getstream.LogLevelInfo
	}
	logger := getstream.NewDefaultLogger(os.Stderr, "", log.LstdFlags, level)
	return logger
}

func validateGlobalArgs(globalArgs *GlobalArgs, command string) error {
	if globalArgs.InputFile == "" && globalArgs.InputDir == "" && globalArgs.InputS3 == "" {
		return fmt.Errorf("either --inputFile or --inputDir or --inputS3 must be specified")
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
		return fmt.Errorf("--inputFile, --inputDir and --inputS3 are exclusive, only one is allowed")
	}

	// --output is optional for list-tracks command (it only displays information)
	if command != "list-tracks" && globalArgs.Output == "" {
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
		return nil, fmt.Errorf("only one filter can be specified at a time: --userId, --sessionId, and --trackId are mutually exclusive")
	}

	var inputPath string
	if globalArgs.InputFile != "" {
		inputPath = globalArgs.InputFile
	} else if globalArgs.InputDir != "" {
		inputPath = globalArgs.InputDir
	} else {
		// TODO: Handle S3 validation
		return nil, fmt.Errorf("Not implemented for now")
	}

	// Parse metadata to validate the single specified argument
	logger := setupLogger(false) // Use non-verbose for validation
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
			return nil, fmt.Errorf("trackID '%s' not found in recording. Use 'list-tracks --format tracks' to see available track IDs", trackID)
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
			return nil, fmt.Errorf("sessionID '%s' not found in recording. Use 'list-tracks --format sessions' to see available session IDs", sessionID)
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
			return nil, fmt.Errorf("userID '%s' not found in recording. Use 'list-tracks --format users' to see available user IDs", userID)
		}
	}

	return metadata, nil
}

func printGlobalUsage() {
	fmt.Fprintf(os.Stderr, "Raw Recording Post Processing Tools\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [global options] <command> [command options]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Global Options:\n")
	fmt.Fprintf(os.Stderr, "  --inputFile <path>     Specify raw recording zip file on file system\n")
	fmt.Fprintf(os.Stderr, "  --inputS3 <path>       Specify raw recording zip file on S3\n")
	fmt.Fprintf(os.Stderr, "  --output <dir>         Specify an output directory (optional for list-tracks)\n")
	fmt.Fprintf(os.Stderr, "  --verbose              Enable verbose logging\n\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  list-tracks            Return list of userId - sessionId - trackId - trackType\n")
	fmt.Fprintf(os.Stderr, "  extract-audio          Generate a playable audio file (webm, mp3, ...)\n")
	fmt.Fprintf(os.Stderr, "  extract-video          Generate a playable video file (webm, mp4, ...)\n")
	fmt.Fprintf(os.Stderr, "  mux-av                 Mux audio and video tracks\n")
	fmt.Fprintf(os.Stderr, "  mix-audio              Mix multiple audio tracks into one file (supports mutually exclusive filters)\n")
	fmt.Fprintf(os.Stderr, "  process-all            Process audio, video, and mux (all-in-one)\n")
	fmt.Fprintf(os.Stderr, "  completion             Generate shell completion scripts\n")
	fmt.Fprintf(os.Stderr, "  help                   Show this help message\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  %s --inputFile recording.zip list-tracks\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s --inputFile recording.zip --output ./out extract-audio --userId user123\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s --inputFile recording.zip --output ./out mix-audio --sessionId session456\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s --verbose --inputFile recording.zip --output ./out mux-av --trackId track789\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Use '%s [global options] <command> --help' for command-specific options.\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nCompletion Setup:\n")
	fmt.Fprintf(os.Stderr, "  # Bash\n")
	fmt.Fprintf(os.Stderr, "  source <(%s completion bash)\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  # Zsh\n")
	fmt.Fprintf(os.Stderr, "  source <(%s completion zsh)\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  # Fish\n")
	fmt.Fprintf(os.Stderr, "  %s completion fish | source\n", os.Args[0])
}

func printHelpIfAsked(args []string, fn func()) {
	// Check for help flag before parsing
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			fn()
			os.Exit(0)
		}
	}
}
func runCompletion(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: raw-tools completion <shell>\n")
		fmt.Fprintf(os.Stderr, "Supported shells: bash, zsh, fish\n")
		os.Exit(1)
	}

	shell := args[0]
	generateCompletion(shell)
}
