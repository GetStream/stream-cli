#! /bin/bash

stream chat:channel:create --channel=$(openssl rand -hex 12) --type="messaging" --name="CLI" --json | jq '.'
