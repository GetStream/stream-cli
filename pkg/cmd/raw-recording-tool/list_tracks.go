package rawrecording

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool/processing"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func listTracksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-tracks",
		Short: "List all tracks in the raw recording with their metadata",
		Long: heredoc.Doc(`
			List all tracks in the raw recording with their metadata.

			This command displays information about all audio and video tracks
			in the recording, including user IDs, session IDs, track IDs, and codecs.

			Note: --output is optional for this command (only displays information).
		`),
		Example: heredoc.Doc(`
			# List all tracks in table format
			$ stream-cli video raw-recording list-tracks --input-file recording.zip

			# Get JSON output for programmatic use
			$ stream-cli video raw-recording list-tracks --input-file recording.zip --format json

			# Get user IDs only
			$ stream-cli video raw-recording list-tracks --input-file recording.zip --format users

			# Filter by track type
			$ stream-cli video raw-recording list-tracks --input-file recording.zip --track-type audio
		`),
		RunE: runListTracks,
	}

	fl := cmd.Flags()
	fl.String(FlagFormat, DefaultFormat, DescFormat)
	fl.String(FlagTrackType, "", DescTrackType)
	fl.String(FlagCompletionType, DefaultCompletionType, DescCompletionType)

	return cmd
}

func runListTracks(cmd *cobra.Command, args []string) error {
	globalArgs, err := getGlobalArgs(cmd)
	if err != nil {
		return err
	}

	// Validate global args (output is optional for list-tracks)
	if err := validateGlobalArgs(globalArgs, false); err != nil {
		return err
	}

	format, _ := cmd.Flags().GetString(FlagFormat)
	trackType, _ := cmd.Flags().GetString(FlagTrackType)
	completionType, _ := cmd.Flags().GetString(FlagCompletionType)

	logger := setupLogger(globalArgs.Verbose)
	logger.Info("Starting list-tracks command")

	// Parse the recording metadata using efficient metadata-only approach
	var inputPath string
	if globalArgs.InputFile != "" {
		inputPath = globalArgs.InputFile
	} else if globalArgs.InputDir != "" {
		inputPath = globalArgs.InputDir
	} else {
		return fmt.Errorf("S3 input not implemented yet")
	}

	parser := processing.NewMetadataParser(logger)
	metadata, err := parser.ParseMetadataOnly(inputPath)
	if err != nil {
		return fmt.Errorf("failed to parse recording: %w", err)
	}

	// Filter tracks if track type is specified
	tracks := processing.FilterTracks(metadata.Tracks, "", "", "", trackType, "")

	// Output in requested format
	switch format {
	case "table":
		printTracksTable(cmd, tracks)
	case "json":
		printTracksJSON(cmd, metadata)
	case "completion":
		printCompletion(cmd, metadata, completionType)
	case "users":
		printUsers(cmd, metadata.UserIDs)
	case "sessions":
		printSessions(cmd, metadata.Sessions)
	case "tracks":
		printTrackIDs(cmd, tracks)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}

	logger.Info("List tracks command completed")
	return nil
}

// printTracksTable prints tracks in a human-readable table format
func printTracksTable(cmd *cobra.Command, tracks []*processing.TrackInfo) {
	if len(tracks) == 0 {
		cmd.Println("No tracks found.")
		return
	}

	// Print header
	cmd.Printf("%-22s %-38s %-38s %-6s %-12s %-15s %-8s\n", "USER ID", "SESSION ID", "TRACK ID", "TYPE", "SCREENSHARE", "CODEC", "SEGMENTS")
	cmd.Printf("%-22s %-38s %-38s %-6s %-12s %-15s %-8s\n",
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
		cmd.Printf("%-22s %-38s %-38s %-6s %-12s %-15s %-8d\n",
			truncateString(track.UserID, 22),
			truncateString(track.SessionID, 38),
			truncateString(track.TrackID, 38),
			track.TrackType,
			screenshareStatus,
			track.Codec,
			track.SegmentCount)
	}
}

// truncateString truncates a string to a maximum length, adding "..." if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// printTracksJSON prints the full metadata in JSON format
func printTracksJSON(cmd *cobra.Command, metadata *processing.RecordingMetadata) {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		cmd.PrintErrf("Error marshaling JSON: %v\n", err)
		return
	}
	cmd.Println(string(data))
}

// printCompletion prints completion-friendly output
func printCompletion(cmd *cobra.Command, metadata *processing.RecordingMetadata, completionType string) {
	switch completionType {
	case "users":
		printUsers(cmd, metadata.UserIDs)
	case "sessions":
		printSessions(cmd, metadata.Sessions)
	case "tracks":
		printTrackIDs(cmd, metadata.Tracks)
	default:
		cmd.PrintErrf("Unknown completion type: %s\n", completionType)
	}
}

// printUsers prints user IDs, one per line
func printUsers(cmd *cobra.Command, userIDs []string) {
	sort.Strings(userIDs)
	for _, userID := range userIDs {
		cmd.Println(userID)
	}
}

// printSessions prints session IDs, one per line
func printSessions(cmd *cobra.Command, sessions []string) {
	sort.Strings(sessions)
	for _, session := range sessions {
		cmd.Println(session)
	}
}

// printTrackIDs prints unique track IDs, one per line
func printTrackIDs(cmd *cobra.Command, tracks []*processing.TrackInfo) {
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
		cmd.Println(trackID)
	}
}
