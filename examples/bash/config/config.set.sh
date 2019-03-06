#! /bin/bash

stream config:set --name="Nick Parsons" --email="nick@getstream.io" --key="foo" --secret="bar" --json | jq '.'
