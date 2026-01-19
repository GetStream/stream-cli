package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/GetStream/getstream-go/v3"
	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
)

type ListTracksArgs struct {
	Format         string // "table", "json", "completion", "users", "sessions", "tracks"
	TrackType      string // Filter by track type: "audio", "video", or "" for all
	CompletionType string // For completion format: "users", "sessions", "tracks"
}

type ListTracksProcess struct {
	logger *getstream.DefaultLogger
}

func NewListTracksProcess(logger *getstream.DefaultLogger) *ListTracksProcess {
	return &ListTracksProcess{logger: logger}
}

func (p *ListTracksProcess) runListTracks(args []string, globalArgs *GlobalArgs) {
	printHelpIfAsked(args, p.printUsage)

	// Parse command-specific flags
	fs := flag.NewFlagSet("list-tracks", flag.ExitOnError)
	listTracksArgs := &ListTracksArgs{}
	fs.StringVar(&listTracksArgs.Format, "format", "table", "Output format: table, json, completion, users, sessions, tracks")
	fs.StringVar(&listTracksArgs.TrackType, "trackType", "", "Filter by track type: audio, video")
	fs.StringVar(&listTracksArgs.CompletionType, "completionType", "tracks", "For completion format: users, sessions, tracks")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	logger := setupLogger(globalArgs.Verbose)

	logger.Info("Starting list-tracks command")

	// Parse the recording metadata using efficient metadata-only approach
	var inputPath string
	if globalArgs.InputFile != "" {
		inputPath = globalArgs.InputFile
	} else if globalArgs.InputDir != "" {
		inputPath = globalArgs.InputDir
	} else {
		// TODO: Handle S3 input
		return // For now, only support local files
	}

	// Use efficient metadata-only parsing (optimized for list-tracks)
	parser := processing.NewMetadataParser(logger)
	metadata, err := parser.ParseMetadataOnly(inputPath)
	if err != nil {
		logger.Error("Failed to parse recording: %v", err)
	}

	// Filter tracks if track type is specified
	tracks := processing.FilterTracks(metadata.Tracks, "", "", "", listTracksArgs.TrackType, "")

	// Output in requested format
	switch listTracksArgs.Format {
	case "table":
		p.printTracksTable(tracks)
	case "json":
		p.printTracksJSON(metadata)
	case "completion":
		p.printCompletion(metadata, listTracksArgs.CompletionType)
	case "users":
		p.printUsers(metadata.UserIDs)
	case "sessions":
		p.printSessions(metadata.Sessions)
	case "tracks":
		p.printTrackIDs(tracks)
	default:
		fmt.Fprintf(os.Stderr, "Unknown format: %s\n", listTracksArgs.Format)
		os.Exit(1)
	}

	logger.Info("List tracks command completed")
}

// printTracksTable prints tracks in a human-readable table format
func (p *ListTracksProcess) printTracksTable(tracks []*processing.TrackInfo) {
	if len(tracks) == 0 {
		fmt.Println("No tracks found.")
		return
	}

	// Print header
	fmt.Printf("%-22s %-38s %-38s %-6s %-12s %-15s %-8s\n", "USER ID", "SESSION ID", "TRACK ID", "TYPE", "SCREENSHARE", "CODEC", "SEGMENTS")
	fmt.Printf("%-22s %-38s %-38s %-6s %-12s %-15s %-8s\n",
		strings.Repeat("-", 22),
		strings.Repeat("-", 38),
		strings.Repeat("-", 38),
		strings.Repeat("-", 6),
		strings.Repeat("-", 12),
		strings.Repeat("-", 15),
		strings.Repeat("-", 8))

	// Print tracks
	for _, track := range tracks {
		screenshareStatus := "No"
		if track.IsScreenshare {
			screenshareStatus = "Yes"
		}
		fmt.Printf("%-22s %-38s %-38s %-6s %-12s %-15s %-8d\n",
			p.truncateString(track.UserID, 22),
			p.truncateString(track.SessionID, 38),
			p.truncateString(track.TrackID, 38),
			track.TrackType,
			screenshareStatus,
			track.Codec,
			track.SegmentCount)
	}
}

// truncateString truncates a string to a maximum length, adding "..." if needed
func (p *ListTracksProcess) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// printTracksJSON prints the full metadata in JSON format
func (p *ListTracksProcess) printTracksJSON(metadata *processing.RecordingMetadata) {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// printCompletion prints completion-friendly output
func (p *ListTracksProcess) printCompletion(metadata *processing.RecordingMetadata, completionType string) {
	switch completionType {
	case "users":
		p.printUsers(metadata.UserIDs)
	case "sessions":
		p.printSessions(metadata.Sessions)
	case "tracks":
		trackIDs := make([]string, 0)
		for _, track := range metadata.Tracks {
			trackIDs = append(trackIDs, track.TrackID)
		}
		// Remove duplicates and sort
		uniqueTrackIDs := p.removeDuplicates(trackIDs)
		sort.Strings(uniqueTrackIDs)
		p.printTrackIDs(metadata.Tracks)
	default:
		fmt.Fprintf(os.Stderr, "Unknown completion type: %s\n", completionType)
	}
}

// printUsers prints user IDs, one per line
func (p *ListTracksProcess) printUsers(userIDs []string) {
	sort.Strings(userIDs)
	for _, userID := range userIDs {
		fmt.Println(userID)
	}
}

// printSessions prints session IDs, one per line
func (p *ListTracksProcess) printSessions(sessions []string) {
	sort.Strings(sessions)
	for _, session := range sessions {
		fmt.Println(session)
	}
}

// printTrackIDs prints unique track IDs, one per line
func (p *ListTracksProcess) printTrackIDs(tracks []*processing.TrackInfo) {
	trackIDs := make([]string, 0)
	seen := make(map[string]bool)

	for _, track := range tracks {
		if !seen[track.TrackID] {
			trackIDs = append(trackIDs, track.TrackID)
			seen[track.TrackID] = true
		}
	}

	sort.Strings(trackIDs)
	for _, trackID := range trackIDs {
		fmt.Println(trackID)
	}
}

// removeDuplicates removes duplicate strings from a slice
func (p *ListTracksProcess) removeDuplicates(input []string) []string {
	keys := make(map[string]bool)
	result := make([]string, 0)

	for _, item := range input {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

func (p *ListTracksProcess) printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: raw-tools [global options] list-tracks [command options]\n\n")
	fmt.Fprintf(os.Stderr, "List all tracks in the raw recording with their metadata.\n")
	fmt.Fprintf(os.Stderr, "Note: --output is optional for this command (only displays information).\n\n")
	fmt.Fprintf(os.Stderr, "Command Options:\n")
	fmt.Fprintf(os.Stderr, "  --format <format>      Output format (default: table)\n")
	fmt.Fprintf(os.Stderr, "                         table    - Human readable table\n")
	fmt.Fprintf(os.Stderr, "                         json     - JSON format\n")
	fmt.Fprintf(os.Stderr, "                         users    - List of user IDs only\n")
	fmt.Fprintf(os.Stderr, "                         sessions - List of session IDs only\n")
	fmt.Fprintf(os.Stderr, "                         tracks   - List of track IDs only\n")
	fmt.Fprintf(os.Stderr, "                         completion - Shell completion format\n")
	fmt.Fprintf(os.Stderr, "  --trackType <type>     Filter by track type: audio, video\n")
	fmt.Fprintf(os.Stderr, "  --completionType <type> For completion format: users, sessions, tracks\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  # List all tracks in table format (no output directory needed)\n")
	fmt.Fprintf(os.Stderr, "  raw-tools --inputFile recording.zip list-tracks\n\n")
	fmt.Fprintf(os.Stderr, "  # Get JSON output for programmatic use\n")
	fmt.Fprintf(os.Stderr, "  raw-tools --inputFile recording.zip list-tracks --format json\n\n")
	fmt.Fprintf(os.Stderr, "  # Get user IDs for completion\n")
	fmt.Fprintf(os.Stderr, "  raw-tools --inputFile recording.zip list-tracks --format users\n")
	fmt.Fprintf(os.Stderr, "\nGlobal Options: Use 'raw-tools --help' to see global options.\n")
}
