#!/usr/bin/env bash

go build ./cmd/stream-cli

random_chars=$(xxd -l16 -ps /dev/urandom)

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

printf "\n\n   #### List configs ####\n\n"
./stream-cli config list
