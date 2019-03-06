#! /bin/bash

stream chat:channel:create --channel=$(openssl rand -hex 32) --type="messaging" --name="CLI" --json
