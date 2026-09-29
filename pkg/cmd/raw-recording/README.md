# Raw Recording CLI

Post-processing tools for Stream Video raw call recordings. Extract, process, and combine audio/video tracks from raw recording archives.

## Features

- **Discovery**: Use `list-tracks` to explore recording contents with screenshare detection
- **Smart Completion**: Shell completion with dynamic values from actual recordings
- **Validation**: Automatic validation of user inputs against available data
- **Multiple Formats**: Support for different output formats (table, JSON, completion)
- **Advanced Processing**: Extract, mux, mix and process audio/video with gap filling
- **S3 Support**: Download recordings directly from S3 or presigned URLs with caching

## Commands

### `list-tracks` - Discovery & Exploration

The `list-tracks` command shows all tracks in a recording with their metadata.

```bash
# List all tracks in table format
stream-cli video raw-recording list-tracks --input-file recording.tar.gz

# Get JSON output for programmatic use
stream-cli video raw-recording list-tracks --input-file recording.tar.gz --format json

# Get user IDs only
stream-cli video raw-recording list-tracks --input-file recording.tar.gz --format users

# Get session IDs only
stream-cli video raw-recording list-tracks --input-file recording.tar.gz --format sessions

# Get track IDs only
stream-cli video raw-recording list-tracks --input-file recording.tar.gz --format tracks

# Filter by track type
stream-cli video raw-recording list-tracks --input-file recording.tar.gz --track-type audio
```

**Options:**
- `--format <format>` - Output format: `table` (default), `json`, `users`, `sessions`, `tracks`, `completion`
- `--track-type <type>` - Filter by track type: `audio`, `video`

**Output Formats:**
- `table` - Human-readable table with screenshare detection (default)
- `json` - Full metadata in JSON format for scripting
- `users` - List of user IDs only (for shell scripts)
- `sessions` - List of session IDs only (for automation)
- `tracks` - List of track IDs only (for filtering)

### `extract-audio` - Extract Audio Tracks

Extract and convert audio tracks from raw recordings to playable MKV format.

```bash
# Extract audio for all users
stream-cli video raw-recording extract-audio --input-file recording.tar.gz --output ./out

# Extract audio for specific user
stream-cli video raw-recording extract-audio --input-file recording.tar.gz --output ./out --user-id user123

# Extract audio for specific session
stream-cli video raw-recording extract-audio --input-file recording.tar.gz --output ./out --session-id session456

# Extract a specific track
stream-cli video raw-recording extract-audio --input-file recording.tar.gz --output ./out --track-id track789

# Disable gap filling
stream-cli video raw-recording extract-audio --input-file recording.tar.gz --output ./out --fill-gaps=false
```

**Options:**
- `--user-id <id>` - Filter by user ID (all tracks for that user)
- `--session-id <id>` - Filter by session ID (all tracks for that session)
- `--track-id <id>` - Filter by track ID (specific track only)
- `--fill-gaps` - Fill temporal gaps with silence when track was muted (default: true)
- `--fix-dtx` - Fix DTX (Discontinuous Transmission) shrink audio (default: true)

**Note**: Filters are mutually exclusive - only one of `--user-id`, `--session-id`, or `--track-id` can be specified at a time.

### `extract-video` - Extract Video Tracks

Extract and convert video tracks from raw recordings to playable MKV format.

```bash
# Extract video for all users
stream-cli video raw-recording extract-video --input-file recording.tar.gz --output ./out

# Extract video for specific user
stream-cli video raw-recording extract-video --input-file recording.tar.gz --output ./out --user-id user123

# Extract video for specific session
stream-cli video raw-recording extract-video --input-file recording.tar.gz --output ./out --session-id session456

# Extract a specific track
stream-cli video raw-recording extract-video --input-file recording.tar.gz --output ./out --track-id track789

# Disable gap filling
stream-cli video raw-recording extract-video --input-file recording.tar.gz --output ./out --fill-gaps=false
```

**Options:**
- `--user-id <id>` - Filter by user ID (all tracks for that user)
- `--session-id <id>` - Filter by session ID (all tracks for that session)
- `--track-id <id>` - Filter by track ID (specific track only)
- `--fill-gaps` - Fill temporal gaps with black frames when track was muted (default: true)

**Note**: Filters are mutually exclusive - only one of `--user-id`, `--session-id`, or `--track-id` can be specified at a time.

### `mux-av` - Combine Audio and Video

Combine audio and video tracks into synchronized files.

```bash
# Mux all tracks
stream-cli video raw-recording mux-av --input-file recording.tar.gz --output ./out

# Mux tracks for specific user
stream-cli video raw-recording mux-av --input-file recording.tar.gz --output ./out --user-id user123

# Mux only user camera tracks (not screenshare)
stream-cli video raw-recording mux-av --input-file recording.tar.gz --output ./out --media user

# Mux only display/screenshare tracks
stream-cli video raw-recording mux-av --input-file recording.tar.gz --output ./out --media display
```

**Options:**
- `--user-id <id>` - Filter by user ID
- `--session-id <id>` - Filter by session ID
- `--track-id <id>` - Filter by track ID
- `--media <type>` - Filter by media type: `user` (camera/microphone), `display` (screenshare), or `both` (default)

**Note**: Filters are mutually exclusive.

### `mix-audio` - Mix Multiple Audio Tracks

Mix audio from multiple users/sessions into a single synchronized audio file.

```bash
# Mix all audio tracks from all users
stream-cli video raw-recording mix-audio --input-file recording.tar.gz --output ./out

# Mix with verbose logging
stream-cli video raw-recording mix-audio --input-file recording.tar.gz --output ./out --verbose
```

Creates `composite_{callType}_{callId}_audio_{timestamp}.mkv` with all tracks properly synchronized based on original timing.

### `process-all` - Complete Workflow

Execute audio extraction, video extraction, and muxing in a single command.

```bash
# Process all tracks
stream-cli video raw-recording process-all --input-file recording.tar.gz --output ./out

# Process tracks for specific user
stream-cli video raw-recording process-all --input-file recording.tar.gz --output ./out --user-id user123

# Process tracks for specific session
stream-cli video raw-recording process-all --input-file recording.tar.gz --output ./out --session-id session456
```

**Options:**
- `--user-id <id>` - Filter by user ID
- `--session-id <id>` - Filter by session ID
- `--track-id <id>` - Filter by track ID

**Output files:**
- `individual_{callType}_{callId}_{userId}_{sessionId}_audio_only_{timestamp}.mkv` - Audio-only files
- `individual_{callType}_{callId}_{userId}_{sessionId}_video_only_{timestamp}.mkv` - Video-only files
- `individual_{callType}_{callId}_{userId}_{sessionId}_audio_video_{timestamp}.mkv` - Combined audio+video files
- `composite_{callType}_{callId}_audio_{timestamp}.mkv` - Mixed audio from all participants

## Global Options

These options are available for all commands:

- `--input-file <path>` - Path to raw recording tar.gz archive
- `--input-dir <path>` - Path to extracted raw recording directory
- `--input-s3 <url>` - S3 URL (`s3://bucket/path`) or presigned HTTPS URL
- `--output <path>` - Output directory (required for most commands)
- `--verbose` - Enable verbose logging
- `--cache-dir <path>` - Cache directory for S3 downloads

**Input options**: Only one of `--input-file`, `--input-dir`, or `--input-s3` can be specified.

## S3 Support

Download recordings directly from S3:

```bash
# Using S3 URL (requires AWS credentials)
stream-cli video raw-recording list-tracks --input-s3 s3://mybucket/recordings/call.tar.gz

# Using presigned HTTPS URL
stream-cli video raw-recording list-tracks --input-s3 "https://mybucket.s3.amazonaws.com/recordings/call.tar.gz?..."

# S3 downloads are cached locally to avoid re-downloading
stream-cli video raw-recording process-all --input-s3 s3://mybucket/call.tar.gz --output ./out
```

## Workflow Examples

### Extract Audio for Each Participant

```bash
# 1. Discover participants
stream-cli video raw-recording list-tracks --input-file call.tar.gz --format users

# 2. Extract each participant's audio
for user in $(stream-cli video raw-recording list-tracks --input-file call.tar.gz --format users); do
    echo "Extracting audio for user: $user"
    stream-cli video raw-recording extract-audio --input-file call.tar.gz --output ./extracted --user-id "$user"
done
```

### Conference Call Audio Mixing

```bash
# Mix all participants into single audio file
stream-cli video raw-recording mix-audio --input-file conference.tar.gz --output ./mixed

# Create session-by-session mixed audio
for session in $(stream-cli video raw-recording list-tracks --input-file conference.tar.gz --format sessions); do
    stream-cli video raw-recording mix-audio --input-file conference.tar.gz --output "./mixed/$session"
done
```

### Complete Processing Pipeline

```bash
# All-in-one processing
stream-cli video raw-recording process-all --input-file recording.tar.gz --output ./complete

# Results:
# - ./complete/individual_*_audio_only_*.mkv (individual audio tracks)
# - ./complete/individual_*_video_only_*.mkv (individual video tracks)
# - ./complete/individual_*_audio_video_*.mkv (combined A/V tracks)
# - ./complete/composite_*_audio_*.mkv (mixed audio)
```

## Dependencies

### FFmpeg

Required for media processing and conversion. Must be compiled with the following libraries:

- `libopus` - Opus audio codec
- `libvpx` - VP8/VP9 video codecs
- `libx264` - H.264 video codec
- `libaom` - AV1 video codec (libaom-av1)
- `libmp3lame` - MP3 audio codec (optional, for MP3 output)

**macOS:**
```bash
brew install ffmpeg
```

**Ubuntu/Debian:**
```bash
sudo apt install ffmpeg
```

### GStreamer

Required for RTP dump to container conversion. Install GStreamer 1.0 with the following plugin packages:

**macOS:**
```bash
brew install gstreamer gst-plugins-base gst-plugins-good gst-plugins-bad gst-plugins-ugly
```

**Ubuntu/Debian:**
```bash
sudo apt install gstreamer1.0-tools gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly
```

**Required GStreamer plugins:**
- `gst-plugins-base` - Core elements (tcpserversrc, filesink)
- `gst-plugins-good` - RTP plugins (rtpjitterbuffer, rtpvp8depay, rtpvp9depay, rtpopusdepay)
- `gst-plugins-bad` - Additional codecs (rtpav1depay, av1parse, matroskamux)
- `gst-plugins-ugly` - H.264 support (rtph264depay, h264parse)

### AWS Credentials (Optional)

Required for S3 URL support (`s3://...`). Not needed for presigned HTTPS URLs.

Configure via:
- Environment variables: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`
- AWS credentials file: `~/.aws/credentials`
- IAM role (when running on AWS infrastructure)

### Go

Go 1.19+ required for building the CLI tool from source.
