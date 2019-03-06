#! /bin/bash

run=$(stream config:get --json)

name=$(jq --raw-output '.name' <<< "${run}")
email=$(jq --raw-output '.email' <<< "${run}")
apiKey=$(jq --raw-output '.apiKey' <<< "${run}")
apiSecret=$(jq --raw-output '.apiSecret' <<< "${run}")

echo $name
echo $email
echo $apiKey
echo $apiSecret
