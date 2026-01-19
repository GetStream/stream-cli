package main

import (
	"fmt"
	"os"
)

// generateCompletion generates shell completion scripts
func generateCompletion(shell string) {
	switch shell {
	case "bash":
		generateBashCompletion()
	case "zsh":
		generateZshCompletion()
	case "fish":
		generateFishCompletion()
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unsupported shell: %s\n", shell)
		_, _ = fmt.Fprintf(os.Stderr, "Supported shells: bash, zsh, fish\n")
		os.Exit(1)
	}
}

// generateBashCompletion generates bash completion script
func generateBashCompletion() {
	script := `#!/bin/bash

_raw_tools_completion() {
    local cur prev words cword
    _init_completion || return

    # Complete subcommands
    if [[ $cword -eq 1 ]]; then
        COMPREPLY=($(compgen -W "list-tracks extract-audio extract-video mux-av help" -- "$cur"))
        return
    fi

    local cmd="${words[1]}"
    
    case "$prev" in
        --inputFile)
            COMPREPLY=($(compgen -f -X "!*.zip" -- "$cur"))
            return
            ;;
        --output)
            COMPREPLY=($(compgen -d -- "$cur"))
            return
            ;;
        --format)
            case "$cmd" in
                list-tracks)
                    COMPREPLY=($(compgen -W "table json completion users sessions tracks" -- "$cur"))
                    ;;
            esac
            return
            ;;
        --trackType)
            COMPREPLY=($(compgen -W "audio video" -- "$cur"))
            return
            ;;
        --userId|--sessionId|--trackId)
            # Dynamic completion using list-tracks
            if [[ -n "${_RAW_TOOLS_INPUT_FILE:-}" ]]; then
                local completion_type=""
                case "$prev" in
                    --userId) completion_type="users" ;;
                    --sessionId) completion_type="sessions" ;;
                    --trackId) completion_type="tracks" ;;
                esac
                if [[ -n "$completion_type" ]]; then
                    local values=$(raw-tools --inputFile "$_RAW_TOOLS_INPUT_FILE" --output /tmp list-tracks --format "$completion_type" 2>/dev/null)
                    COMPREPLY=($(compgen -W "$values" -- "$cur"))
                fi
            else
                COMPREPLY=()
            fi
            return
            ;;
    esac

    # Complete global flags
    local global_flags="--inputFile --inputS3 --output --verbose --help"
    local cmd_flags=""
    
    case "$cmd" in
        list-tracks)
            cmd_flags="--format --trackType --completionType"
            ;;
        extract-audio|extract-video)
            cmd_flags="--userId --sessionId --trackId --fill_gaps"
            ;;
        mux-av)
            cmd_flags="--userId --sessionId --trackId --media"
            ;;
        mix-audio)
            cmd_flags=""
            ;;
    esac
    
    COMPREPLY=($(compgen -W "$global_flags $cmd_flags" -- "$cur"))
}

# Store input file for dynamic completion
_raw_tools_set_input_file() {
    local i
    for (( i=1; i < ${#COMP_WORDS[@]}; i++ )); do
        if [[ "${COMP_WORDS[i]}" == "--inputFile" && i+1 < ${#COMP_WORDS[@]} ]]; then
            export _RAW_TOOLS_INPUT_FILE="${COMP_WORDS[i+1]}"
            break
        fi
    done
}

# Hook to set input file before completion
complete -F _raw_tools_completion raw-tools

# Wrapper to set input file
_raw_tools_wrapper() {
    _raw_tools_set_input_file
    _raw_tools_completion "$@"
}

complete -F _raw_tools_wrapper raw-tools`

	fmt.Println(script)
}

// generateZshCompletion generates zsh completion script
func generateZshCompletion() {
	script := `#compdef raw-tools

_raw_tools() {
    local context state line
    typeset -A opt_args

    _arguments -C \
        '1: :_raw_tools_commands' \
        '*:: :->args'

    case $state in
        args)
            case $words[1] in
                list-tracks)
                    _raw_tools_list_tracks
                    ;;
                extract-audio|extract-video)
                    _raw_tools_extract
                    ;;
                mux-av)
                    _raw_tools_mux_av
                    ;;
            esac
            ;;
    esac
}

_raw_tools_commands() {
    local commands=(
        'list-tracks:List all tracks with metadata'
        'extract-audio:Generate playable audio files'
        'extract-video:Generate playable video files'
        'mux-av:Mux audio and video tracks'
        'help:Show help'
    )
    _describe 'commands' commands
}

_raw_tools_global_args() {
    _arguments \
        '--inputFile[Specify raw recording zip file]:file:_files -g "*.zip"' \
        '--inputS3[Specify raw recording zip file on S3]:s3path:' \
        '--output[Specify output directory]:directory:_directories' \
        '--verbose[Enable verbose logging]' \
        '--help[Show help]'
}

_raw_tools_list_tracks() {
    _arguments \
        '--format[Output format]:format:(table json completion users sessions tracks)' \
        '--trackType[Filter by track type]:type:(audio video)' \
        '--completionType[Completion type]:type:(users sessions tracks)' \
        '*: :_raw_tools_global_args'
}

_raw_tools_extract() {
    _arguments \
        '--userId[User ID filter]:userid:_raw_tools_complete_users' \
        '--sessionId[Session ID filter]:sessionid:_raw_tools_complete_sessions' \
        '--trackId[Track ID filter]:trackid:_raw_tools_complete_tracks' \
        '--fill_gaps[Fill gaps with silence/black frames]' \
        '*: :_raw_tools_global_args'
}

_raw_tools_mux_av() {
    _arguments \
        '--userId[User ID filter]:userid:_raw_tools_complete_users' \
        '--sessionId[Session ID filter]:sessionid:_raw_tools_complete_sessions' \
        '--trackId[Track ID filter]:trackid:_raw_tools_complete_tracks' \
        '--media[Media type]:media:(user display both)' \
        '*: :_raw_tools_global_args'
}
// no mix-audio specific flags

# Dynamic completion helpers
_raw_tools_complete_users() {
    local input_file
    for ((i=1; i <= $#words; i++)); do
        if [[ $words[i] == "--inputFile" && i+1 <= $#words ]]; then
            input_file=$words[i+1]
            break
        fi
    done
    
    if [[ -n "$input_file" ]]; then
        local users=($(raw-tools --inputFile "$input_file" --output /tmp list-tracks --format users 2>/dev/null))
        _wanted users expl 'user ID' compadd "$@" $users
    else
        _wanted users expl 'user ID' compadd "$@"
    fi
}

_raw_tools_complete_sessions() {
    local input_file
    for ((i=1; i <= $#words; i++)); do
        if [[ $words[i] == "--inputFile" && i+1 <= $#words ]]; then
            input_file=$words[i+1]
            break
        fi
    done
    
    if [[ -n "$input_file" ]]; then
        local sessions=($(raw-tools --inputFile "$input_file" --output /tmp list-tracks --format sessions 2>/dev/null))
        _wanted sessions expl 'session ID' compadd "$@" $sessions
    else
        _wanted sessions expl 'session ID' compadd "$@"
    fi
}

_raw_tools_complete_tracks() {
    local input_file
    for ((i=1; i <= $#words; i++)); do
        if [[ $words[i] == "--inputFile" && i+1 <= $#words ]]; then
            input_file=$words[i+1]
            break
        fi
    done
    
    if [[ -n "$input_file" ]]; then
        local tracks=($(raw-tools --inputFile "$input_file" --output /tmp list-tracks --format tracks 2>/dev/null))
        _wanted tracks expl 'track ID' compadd "$@" $tracks
    else
        _wanted tracks expl 'track ID' compadd "$@"
    fi
}

_raw_tools "$@"`

	fmt.Println(script)
}

// generateFishCompletion generates fish completion script
func generateFishCompletion() {
	script := `# Fish completion for raw-tools

# Complete commands
complete -c raw-tools -f -n '__fish_use_subcommand' -a 'list-tracks' -d 'List all tracks with metadata'
complete -c raw-tools -f -n '__fish_use_subcommand' -a 'extract-audio' -d 'Generate playable audio files'
complete -c raw-tools -f -n '__fish_use_subcommand' -a 'extract-video' -d 'Generate playable video files'
complete -c raw-tools -f -n  '__fish_use_subcommand' -a 'mux-av' -d 'Mux audio and video tracks'
complete -c raw-tools -f -n '__fish_use_subcommand' -a 'help' -d 'Show help'

# Global options
complete -c raw-tools -l inputFile -d 'Specify raw recording zip file' -r -F
complete -c raw-tools -l inputS3 -d 'Specify raw recording zip file on S3' -r
complete -c raw-tools -l output -d 'Specify output directory' -r -a '(__fish_complete_directories)'
complete -c raw-tools -l verbose -d 'Enable verbose logging'
complete -c raw-tools -l help -d 'Show help'

# list-tracks specific options
complete -c raw-tools -n '__fish_seen_subcommand_from list-tracks' -l format -d 'Output format' -r -a 'table json completion users sessions tracks'
complete -c raw-tools -n '__fish_seen_subcommand_from list-tracks' -l trackType -d 'Filter by track type' -r -a 'audio video'
complete -c raw-tools -n '__fish_seen_subcommand_from list-tracks' -l completionType -d 'Completion type' -r -a 'users sessions tracks'

# extract commands specific options
complete -c raw-tools -n '__fish_seen_subcommand_from extract-audio extract-video' -l userId -d 'User ID filter' -r
complete -c raw-tools -n '__fish_seen_subcommand_from extract-audio extract-video' -l sessionId -d 'Session ID filter' -r
complete -c raw-tools -n '__fish_seen_subcommand_from extract-audio extract-video' -l trackId -d 'Track ID filter' -r
complete -c raw-tools -n '__fish_seen_subcommand_from extract-audio extract-video' -l fill_gaps -d 'Fill gaps'

# mux-av specific options
complete -c raw-tools -n '__fish_seen_subcommand_from mux-av' -l userId -d 'User ID filter' -r
complete -c raw-tools -n '__fish_seen_subcommand_from mux-av' -l sessionId -d 'Session ID filter' -r
complete -c raw-tools -n '__fish_seen_subcommand_from mux-av' -l trackId -d 'Track ID filter' -r
complete -c raw-tools -n '__fish_seen_subcommand_from mux-av' -l media -d 'Media type' -r -a 'user display both'

# mix-audio has no command-specific options`

	fmt.Println(script)
}
