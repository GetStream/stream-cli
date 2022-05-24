#!/bin/bash

test_curl_exists() {
  if ! curl --version > /dev/null 2>&1; then
    echo "curl is required (it's used to download the CLI). Exiting."
    exit 1
  fi
}

is_macos() {
    [[ "$OSTYPE" == "darwin"* ]]
}

get_platform() {
    linuxOrMac="$(uname -s)"
    intelOrArm="$(uname -m)"
    echo "${linuxOrMac}_${intelOrArm}"
}

# We need curl to download from GitHub Releases
test_curl_exists

# Get the current platform. Such as Linux_x86 or Darwin_arm64
PLATFORM=$(get_platform)

# Grab the download url for that platform
URL=$(curl -s https://api.github.com/repos/GetStream/stream-cli/releases/latest | grep "$PLATFORM"  | cut -d '"' -f 4 | sed '1d')

# Create a folder for the CLI in ~/.stream-cli
echo " > Ensuring folder exists: ~/.stream-cli"
mkdir -p ~/.stream-cli

# Download the CLI
echo " > Downloading the tarball from GitHub."
curl -s -L "$URL" -o ~/.stream-cli/stream-cli.tar.gz

# Extract the tar then remove it
echo " > Extracting the tarball to ~/.stream-cli"
tar -xzf ~/.stream-cli/stream-cli.tar.gz -C ~/.stream-cli
rm ~/.stream-cli/stream-cli.tar.gz

# If MacOS, we need to trust the binary
if is_macos; then
    echo " > In MacOS, the binary needs to be trusted by the system. We use 'xattr -d com.apple.quarantine ~/.stream-cli' command to do that."
    echo " > Do you agree to continue? [y/n]"
    read -r AGREE
    if [ "$AGREE" == "y" ]; then
        xattr -d com.apple.quarantine ~/.stream-cli
    fi 
fi

echo " > The CLI has been installed in ~/.stream-cli."
# shellcheck disable=SC2016
echo ' > As a last step, add it to your $PATH or create a symbolic link to /usr/local/bin/stream-cli.'
echo " > Do you want to create a symbolic link now? (requires sudo) [y/n]"

read -r AGREE
if [ "$AGREE" != "y" ]; then
    echo " > All done. 🎉 You'll find the CLI at ~/.stream-cli/stream-cli."
    exit 0
fi

sudo ln -f -s ~/.stream-cli/stream-cli /usr/local/bin/stream-cli

echo " > All done! 🎉 Try 'stream-cli --help'"
