#!/usr/bin/env bash

set -o pipefail

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
./stream-cli chat list-channels --type messaging --limit 1

printf '\n\n   #### Create channel type ####\n\n'
./stream-cli chat create-channel-type -p "{\"name\": \"$random_chars\"}"

printf '\n\n   #### Get channel type ####\n\n'
./stream-cli chat get-channel-type "$random_chars"

printf '\n\n   #### Update channel type   ####\n\n'
./stream-cli chat update-channel-type --type "$random_chars" --properties "{\"quotes\":true}"

printf '\n\n   #### List channel type ####\n\n'
./stream-cli chat list-channel-types

# Let's make sure this is the last command so we clean up after ourselves
printf '\n\n   #### Delete channel type ####\n\n'
./stream-cli chat delete-channel-type "$random_chars"

printf '\n\n   #### Upsert user  ####\n\n'
./stream-cli chat upsert-user --properties "{\"id\":\"$random_chars\"}"

printf '\n\n   #### Create token for user  ####\n\n'
./stream-cli chat create-token --user "$random_chars"

printf '\n\n   #### Query user  ####\n\n'
./stream-cli chat query-users --filter "{\"id\":\"$random_chars\"}"

printf '\n\n   #### Send message  ####\n\n'
./stream-cli chat send-message -i "$random_chars" -t messaging --text text -u "$random_chars"

printf '\n\n   #### Get a single message  ####\n\n'
./stream-cli chat get-message "$random_chars"

printf '\n\n   #### Cleanup process  ####\n\n'
printf "\n\n   #### Delete channel ####\n\n"
./stream-cli chat delete-channel --type messaging --id "$random_chars" --hard

printf '\n\n   #### Delete user  ####\n\n'
./stream-cli chat delete-user --user "$random_chars" --hard-delete --mark-messages-deleted --delete-conversations
