#!/usr/bin/env bash

go build ./cmd/stream-cli

random_chars=$(xxd -l16 -ps /dev/urandom)

printf "\n\n   #### List configs ####\n\n"
./stream-cli config list

printf "\n\n   #### Update app settings ####\n\n"
./stream-cli chat update-app -p '{"multi_tenant_enabled":true}'

printf "\n\n   #### Get app settings ####\n\n"
./stream-cli chat get-app

printf "\n\n   #### Create channel ####\n\n"
./stream-cli chat create-channel --type messaging --id "$random_chars" --user user

printf "\n\n   #### Get channel ####\n\n"
./stream-cli chat get-channel --type messaging --id "$random_chars"

printf "\n\n   #### Update channel ####\n\n"
./stream-cli chat update-channel --type messaging --id "$random_chars" --properties "{\"frozen\":false}"

printf "\n\n   #### List channels ####\n\n"
./stream-cli chat list-channels --type messaging --limit 5

# Let's make sure this is the last command so we clean up after ourselves
printf "\n\n   #### Delete channel ####\n\n"
./stream-cli chat delete-channel --type messaging --id "$random_chars" --hard

printf '\n\n   #### Create channel type ####\n\n'
./stream-cli chat create-channel-type -p "{\"name\": \"$random_chars\"}"

printf '\n\n   #### Get channel type ####\n\n'
./stream-cli chat get-channel-type --channel-type "$random_chars"

printf '\n\n   #### Update channel type   ####\n\n'
./stream-cli chat update-channel-type --channel-type "$random_chars" --properties "{\"quotes\":true}"

printf '\n\n   #### List channel type ####\n\n'
./stream-cli chat list-channel-types

# Let's make sure this is the last command so we clean up after ourselves
printf '\n\n   #### Delete channel type ####\n\n'
./stream-cli chat delete-channel-type --channel-type "$random_chars"
