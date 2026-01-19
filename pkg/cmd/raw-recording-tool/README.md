# Raw-Tools CLI

Post-processing tools for raw video call recordings with intelligent completion, validation, and advanced audio/video processing.

## Features

- **Discovery**: Use `list-tracks` to explore recording contents with screenshare detection
- **Smart Completion**: Shell completion with dynamic values from actual recordings
- **Validation**: Automatic validation of user inputs against available data
- **Multiple Formats**: Support for different output formats (table, JSON, completion)
- **Advanced Processing**: Extract, mux, mix and process audio/video with gap filling
- **Hybrid Architecture**: Optimized performance for different use cases

## Commands

### `list-tracks` - Discovery & Completion Hub

The `list-tracks` command serves as both a discovery tool and completion engine for other commands.

```bash
# Basic usage - see all tracks in table format (no --output needed)
raw-tools --inputFile recording.tar.gz list-tracks

# Get JSON output for programmatic use
raw-tools --inputFile recording.tar.gz list-tracks --format json

# Get completion-friendly lists
raw-tools --inputFile recording.tar.gz list-tracks --format users
raw-tools --inputFile recording.tar.gz list-tracks --format sessions
raw-tools --inputFile recording.tar.gz list-tracks --format tracks
```

**Options:**
- `--format <format>` - Output format: `table` (default), `json`, `users`, `sessions`, `tracks`, `completion`
- `--trackType <type>` - Filter by track type: `audio`, `video` (optional)
- `-h, --help` - Show help message

**Output Formats:**
- `table` - Human-readable table with screenshare detection (default)
- `json` - Full metadata in JSON format for scripting
- `users` - List of user IDs only (for shell scripts)
- `sessions` - List of session IDs only (for automation)  
- `tracks` - List of track IDs only (for filtering)
- `completion` - Shell completion format

### `extract-audio` - Extract Audio Tracks

Extract and convert individual audio tracks from raw recordings to WebM format.

```bash
# Extract audio for all users
raw-tools --inputFile recording.zip --output ./output extract-audio

# Extract audio for specific user with gap filling
raw-tools --inputFile recording.zip --output ./output extract-audio --userId user123 --fill_gaps

# Extract audio for specific session
raw-tools --inputFile recording.zip --output ./output extract-audio --sessionId session456

# Extract specific track only
raw-tools --inputFile recording.zip --output ./output extract-audio --trackId track789
```

**Options:**
- `--userId <id>` - Filter by user ID (returns all tracks for that user)
- `--sessionId <id>` - Filter by session ID (returns all tracks for that session)
- `--trackId <id>` - Filter by track ID (returns only that specific track)
- **Note**: These filters are mutually exclusive - only one can be specified at a time
- `--fill_gaps` - Fill temporal gaps between segments with silence (recommended for playback)
- `-h, --help` - Show help message

**Mutually Exclusive Filtering:**
- Only one filter can be specified at a time: `--userId`, `--sessionId`, or `--trackId`
- `--trackId` returns exactly one track (the specified track)
- `--sessionId` returns all tracks for that session (multiple tracks possible)  
- `--userId` returns all tracks for that user (multiple tracks possible)
- If no filter is specified, all tracks are processed

### `extract-video` - Extract Video Tracks

Extract and convert individual video tracks from raw recordings to WebM format.

```bash
# Extract video for all users
raw-tools --inputFile recording.zip --output ./output extract-video

# Extract video for specific user with black frame filling
raw-tools --inputFile recording.zip --output ./output extract-video --userId user123 --fill_gaps

# Extract screenshare video only
raw-tools --inputFile recording.zip --output ./output extract-video --userId user456 --fill_gaps
```

**Options:**
- `--userId <id>` - Filter by user ID (returns all tracks for that user)  
- `--sessionId <id>` - Filter by session ID (returns all tracks for that session)
- `--trackId <id>` - Filter by track ID (returns only that specific track)
- **Note**: These filters are mutually exclusive - only one can be specified at a time
- `--fill_gaps` - Fill temporal gaps between segments with black frames (recommended for playback)
- `-h, --help` - Show help message

**Video Processing:**
- Supports regular camera video and screenshare video
- Automatically detects and preserves video codec (VP8, VP9, H264, AV1)
- Gap filling generates black frames matching original video dimensions and framerate

### `mux-av` - Mux Audio/Video

Combine individual audio and video tracks with proper synchronization and timing offsets.

```bash
# Mux audio/video for all users
raw-tools --inputFile recording.zip --output ./output mux-av

# Mux for specific user with proper sync
raw-tools --inputFile recording.zip --output ./output mux-av --userId user123

# Mux for specific session
raw-tools --inputFile recording.zip --output ./output mux-av --sessionId session456

# Mux specific tracks with precise control
raw-tools --inputFile recording.zip --output ./output mux-av --userId user123 --sessionId session456
```

**Options:**
- `--userId <id>` - Filter by user ID (returns all tracks for that user)
- `--sessionId <id>` - Filter by session ID (returns all tracks for that session)
- `--trackId <id>` - Filter by track ID (returns only that specific track)
- **Note**: These filters are mutually exclusive - only one can be specified at a time
- `--media <type>` - Filter by media type: `user` (camera/microphone), `display` (screen sharing), or `both` (default)
- `-h, --help` - Show help message

**Features:**
- Automatic timing synchronization between audio and video using RTCP timestamps
- Gap filling for seamless playback (always enabled for muxing)
- Single combined WebM output per user/session combination
- Intelligent offset calculation for perfect A/V sync
- Supports all video codecs (VP8, VP9, H264, AV1) with Opus audio
- Media type filtering ensures consistent pairing (user camera ↔ user microphone, display sharing ↔ display audio)

**Media Type Examples:**
```bash
# Mux only user camera/microphone tracks
raw-tools --inputFile recording.zip --output ./output mux-av --userId user123 --media user

# Mux only display sharing tracks  
raw-tools --inputFile recording.zip --output ./output mux-av --userId user123 --media display

# Mux both types with proper pairing (default)
raw-tools --inputFile recording.zip --output ./output mux-av --userId user123 --media both
```

### `mix-audio` - Mix Multiple Audio Tracks

Mix audio from multiple users/sessions into a single synchronized audio file, perfect for conference call reconstruction.

```bash
# Mix audio from all users (full conference call)
raw-tools --inputFile recording.zip --output ./output mix-audio

# Mix audio from specific user across all sessions
raw-tools --inputFile recording.zip --output ./output mix-audio --userId user123

# Mix audio from specific session (all users in that session)
raw-tools --inputFile recording.zip --output ./output mix-audio --sessionId session456

# Mix specific tracks with fine control
raw-tools --inputFile recording.zip --output ./output mix-audio --userId user123 --sessionId session456
```

**Options:**
- `--userId <id>` - Filter by user ID (returns all tracks for that user)
- `--sessionId <id>` - Filter by session ID (returns all tracks for that session)
- `--trackId <id>` - Filter by track ID (returns only that specific track)
- **Note**: These filters are mutually exclusive - only one can be specified at a time
- `--no-fill-gaps` - Disable gap filling (not recommended for mixing, gaps enabled by default)
- `-h, --help` - Show help message

**Perfect for:**
- Conference call audio reconstruction with proper timing
- Multi-participant audio analysis and review
- Creating complete session audio timelines
- Audio synchronization testing and validation
- Podcast-style recordings from video calls

**Advanced Mixing:**
- Uses FFmpeg adelay and amix filters for professional-quality mixing
- Automatic timing offset calculation based on segment metadata
- Gap filling with silence maintains temporal relationships
- Output: Single `mixed_audio.webm` file with all tracks properly synchronized

### `process-all` - Complete Workflow

Execute audio extraction, video extraction, and muxing in a single command - the all-in-one solution.

```bash
# Process everything for all users
raw-tools --inputFile recording.zip --output ./output process-all

# Process everything for specific user
raw-tools --inputFile recording.zip --output ./output process-all --userId user123

# Process specific session with all participants
raw-tools --inputFile recording.zip --output ./output process-all --sessionId session456

# Process specific tracks with full workflow
raw-tools --inputFile recording.zip --output ./output process-all --userId user123 --sessionId session456
```

**Options:**
- `--userId <id>` - Filter by user ID (returns all tracks for that user)
- `--sessionId <id>` - Filter by session ID (returns all tracks for that session)
- `--trackId <id>` - Filter by track ID (returns only that specific track)
- **Note**: These filters are mutually exclusive - only one can be specified at a time
- `-h, --help` - Show help message

**Workflow Steps:**
1. **Audio Extraction** - Extracts all matching audio tracks with gap filling enabled
2. **Video Extraction** - Extracts all matching video tracks with gap filling enabled  
3. **Audio/Video Muxing** - Combines corresponding audio and video tracks with sync

**Outputs:**
- Individual audio tracks (WebM format): `audio_userId_sessionId_trackId.webm`
- Individual video tracks (WebM format): `video_userId_sessionId_trackId.webm`
- Combined audio/video files (WebM format): `muxed_userId_sessionId_combined.webm`
- All files include gap filling for seamless playback
- Perfect for bulk processing and automated workflows

## Completion Workflow Architecture

### 1. Discovery Phase
```bash
# First, explore what's in your recording
raw-tools --inputFile recording.zip list-tracks

# Example output with screenshare detection:
# USER ID              SESSION ID           TRACK ID             TYPE    SCREENSHARE  CODEC           SEGMENTS
# -------------------- -------------------- -------------------- ------- ------------ --------------- --------
# user_abc123         session_xyz789       track_001            audio   No           audio/opus      3
# user_abc123         session_xyz789       track_002            video   No           video/VP8       2  
# user_def456         session_xyz789       track_003            video   Yes          video/VP8       1
```

### 2. Shell Completion Setup

```bash
# Install completion for your shell
source <(raw-tools completion bash)   # Bash
source <(raw-tools completion zsh)    # Zsh  
raw-tools completion fish | source    # Fish
```

### 3. Dynamic Completion in Action

With completion enabled, the CLI will:
- **Auto-complete commands** and flags
- **Dynamically suggest user IDs** from the actual recording
- **Validate inputs** against available data
- **Provide helpful error messages** with discovery hints

```bash
# Tab completion will suggest actual user IDs from your recording
raw-tools --inputFile recording.zip --output ./out extract-audio --userId <TAB>
# Shows: user_abc123  user_def456

# Invalid inputs show helpful errors
raw-tools --inputFile recording.zip --output ./out extract-audio --userId invalid_user
# Error: userID 'invalid_user' not found in recording. Available users: user_abc123, user_def456
# Tip: Use 'raw-tools --inputFile recording.zip --output ./out list-tracks --format users' to see available user IDs
```

### 4. Programmatic Integration

```bash
# Get user IDs for scripts
USERS=$(raw-tools --inputFile recording.zip list-tracks --format users)

# Process each user
for user in $USERS; do
    raw-tools --inputFile recording.zip --output ./output extract-audio --userId "$user" --fill_gaps
done

# Get JSON metadata for complex processing
raw-tools --inputFile recording.zip list-tracks --format json > metadata.json
```

## Workflow Examples

### Example 1: Extract Audio for Each Participant

```bash
# 1. Discover participants
raw-tools --inputFile call.zip list-tracks --format users

# 2. Extract each participant's audio
for user in $(raw-tools --inputFile call.zip list-tracks --format users); do
    echo "Extracting audio for user: $user"
    raw-tools --inputFile call.zip --output ./extracted extract-audio --userId "$user" --fill_gaps
done
```

### Example 2: Quality Check Before Processing

```bash
# 1. Get full metadata overview
raw-tools --inputFile recording.zip list-tracks --format json > recording_info.json

# 2. Check track counts
audio_tracks=$(raw-tools --inputFile recording.zip list-tracks --trackType audio --format json | jq '.tracks | length')
video_tracks=$(raw-tools --inputFile recording.zip list-tracks --trackType video --format json | jq '.tracks | length')

echo "Found $audio_tracks audio tracks and $video_tracks video tracks"

# 3. Process only if we have both audio and video
if [ "$audio_tracks" -gt 0 ] && [ "$video_tracks" -gt 0 ]; then
    raw-tools --inputFile recording.zip --output ./output mux-av
fi
```

### Example 3: Conference Call Audio Mixing

```bash
# 1. Mix all participants into single audio file
raw-tools --inputFile conference.zip --output ./mixed mix-audio

# 2. Mix specific users for focused conversation (individual commands)
raw-tools --inputFile conference.zip --output ./mixed mix-audio --userId user1
raw-tools --inputFile conference.zip --output ./mixed mix-audio --userId user2

# 3. Create session-by-session mixed audio
for session in $(raw-tools --inputFile conference.zip list-tracks --format sessions); do
    raw-tools --inputFile conference.zip --output "./mixed/$session" mix-audio --sessionId "$session"
done
```

### Example 4: Complete Processing Pipeline

```bash
# All-in-one processing for the entire recording
raw-tools --inputFile recording.zip --output ./complete process-all

# Results in:
# - ./complete/audio_*.webm (individual audio tracks)
# - ./complete/video_*.webm (individual video tracks)  
# - ./complete/muxed_*.webm (combined A/V tracks)
```

### Example 5: Session-Based Processing

```bash
# 1. Process each session separately
for session in $(raw-tools --inputFile recording.zip list-tracks --format sessions); do
    echo "Processing session: $session"
    
    # Extract all audio from this session
    raw-tools --inputFile recording.zip --output "./output/$session" extract-audio --sessionId "$session" --fill_gaps
    
    # Extract all video from this session  
    raw-tools --inputFile recording.zip --output "./output/$session" extract-video --sessionId "$session" --fill_gaps
    
    # Mux audio/video for this session
    raw-tools --inputFile recording.zip --output "./output/$session" mux-av --sessionId "$session"
done
```

## Architecture & Performance

### Hybrid Processing Architecture

The tool uses an intelligent hybrid approach optimized for different use cases:

**Fast Metadata Reading (`list-tracks`):**
- Direct tar.gz parsing for metadata-only operations
- Skips extraction of large media files (.rtpdump/.sdp)
- 10-50x faster than full extraction for discovery workflows

**Full Processing (extraction commands):**
- Complete archive extraction to temporary directories
- Access to all media files for conversion and processing
- Unified processing pipeline for reliability

### Command Categories

1. **Discovery Commands** (`list-tracks`)
   - Optimized for speed and shell completion
   - Minimal resource usage
   - Instant metadata access

2. **Processing Commands** (`extract-*`, `mix-*`, `mux-*`, `process-all`)
   - Full archive extraction and processing
   - Complete media file access
   - Advanced audio/video operations

3. **Utility Commands** (`completion`, `help`)
   - Shell integration and documentation

## Benefits of the Architecture

1. **Discoverability**: No need to guess user IDs, session IDs, or track IDs
2. **Performance**: Optimized operations for different use cases  
3. **Validation**: Immediate feedback if specified IDs don't exist
4. **Efficiency**: Tab completion speeds up command construction
5. **Reliability**: Prevents typos and invalid commands
6. **Scriptability**: Programmatic access to metadata for automated workflows
7. **User Experience**: Helpful error messages with actionable suggestions
8. **Advanced Processing**: Conference call reconstruction and analysis capabilities

## File Structure

```
cmd/raw-tools/
├── main.go              # Main CLI entry point and routing
├── metadata.go          # Shared metadata parsing and filtering (hybrid architecture)
├── completion.go        # Shell completion scripts generation  
├── list_tracks.go       # Discovery and completion command (optimized)
├── extract_audio.go     # Audio extraction with validation
├── extract_video.go     # Video extraction with validation
├── extract_track.go     # Generic extraction logic (shared)
├── mix_audio.go         # Multi-user audio mixing
├── mux_av.go           # Audio/video synchronization and muxing
├── process_all.go      # All-in-one processing workflow
└── README.md           # This documentation
```

## Dependencies

- **FFmpeg**: Required for media processing and conversion
- **Go 1.19+**: For building the CLI tool
