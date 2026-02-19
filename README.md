![Stream Cli](./assets/logo.png)

# Stream CLI

---
> ## 🚨 **Breaking changes in v1.0 <**
> We have completely rewritten the Node.JS CLI to Go in 2022 Q1. Some of the changes:
> - The installation process is easier since it doesn't have any prerequisites (such as NPM). You can just simply download the executable and run it.
> - The name of the executable is `stream-cli` instead of `stream` to avoid conflicts with an existing tool ([imagemagick](https://github.com/GetStream/stream-cli/issues/33)). But you can rename it if you want to.
> - The command invocation is `stream-cli chat [verb-noun] [args] [options]` instead of `stream [verb:noun] [args] [options]`. The most obvious change is using dash instead of colon. We also added the `chat` keyword to preserve domain for our other product [Feeds](https://getstream.io/activity-feeds/).
> - The 1.0.0 Go version's feature set is matching the old one. But if you miss anything, feel free to open an issue.

Stream's Command Line Interface (CLI) makes it easy to create and manage your [Stream](https://getstream.io) apps directly from the terminal. Currently, only Chat is supported; however, the ability to manage Feeds will be coming soon.

# 📚 Documentation
The full documentation is deployed to [GitHub Pages](https://getstream.github.io/stream-cli/).

# 🗒 Issues

If you're experiencing problems directly related to the CLI, please add an [issue on GitHub](https://github.com/getstream/stream-cli/issues).

For other issues, submit a [support ticket](https://getstream.io/support).

# 📝 Changelog

As with any project, things are always changing. If you're interested in seeing what's changed in the Stream CLI, the changelog for this project can be found [here](./CHANGELOG.md).

# 🏗 Installation

The Stream CLI is written in Go and precompiled into a single binary. It doesn't have any prerequisites.

## Download the binaries
You can find the binaries in the [Release section](https://github.com/GetStream/stream-cli/releases) of this repository. We also wrote a short script to download them and put it to your $PATH.

### Bash (MacOS and Linux) <!-- omit in toc -->
```shell
$ /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/GetStream/stream-cli/master/install/install.sh)"
```

### PowerShell (Windows) <!-- omit in toc -->
```powershell
$ Invoke-WebRequest -Uri "https://raw.githubusercontent.com/GetStream/stream-cli/master/install/install.ps1" -OutFile "install.ps1"; powershell.exe -ExecutionPolicy Bypass -File ./install.ps1
```

## Homebrew

For MacOS users, it's also available via Homebrew:

```shell
$ brew tap GetStream/stream-cli https://github.com/GetStream/stream-cli
$ brew install --cask stream-cli
```

## Compile yourself
```shell
$ git clone git@github.com:GetStream/stream-cli.git
$ cd stream-cli
$ go build ./cmd/stream-cli
$ ./stream-cli --version
stream-cli version 1.0.0
```

# 🚀 Getting Started

In order to initialize the CLI, it's as simple as:

![Stream](./assets/first_config.svg)

> Note: Your API key and secret can be found on the [Stream Dashboard](https://getstream.io/dashboard) and is specific to your organization.

# 🚨 Warning

We purposefully chose the executable name `stream-cli` to avoid conflict with another tool called [`imagemagick`](https://imagemagick.org/index.php) which [already has a `stream` executable](https://github.com/GetStream/stream-cli/issues/33). 

If you do not have `imagemagick` installed, it might be more comfortable to rename `stream-cli` to `stream`. Alternatively you can set up a symbolic link:

```shell
$ ln -s ~/Downloads/stream-cli /usr/local/bin/stream
$ stream --version
stream-cli version 1.0.0
```

# 🔨 Syntax

Basic commands use the following syntax:

```shell
$ stream-cli [chat|feeds] [command] [args] [options]
```

Example:

```shell
$ stream-cli chat get-channel -t messaging -i redteam
```

The `--help` keyword is available every step of the way. Examples:

```shell
$ stream-cli --help
$ stream-cli chat --help
$ stream-cli chat get-channel --help
```

# 💬 Auto completion
We provide autocompletion for the most popular shells (PowerShell, Bash, ZSH, Fish).

```shell
$ stream-cli completion --help
```

# 📣 Feedback

If you have any suggestions or just want to let us know what you think of the Stream CLI, please send us a message at support@getstream.io or create a [GitHub Issue](https://github.com/getstream/stream-cli/issues).

# 🔧 Development

We welcome code changes that improve this library or fix a problem, please make sure to follow all best practices and add tests if applicable before submitting a Pull Request on Github. We are very happy to merge your code in the official repository. Make sure to sign our [Contributor License Agreement (CLA)](https://docs.google.com/forms/d/e/1FAIpQLScFKsKkAJI7mhCr7K9rEIOpqIDThrWxuvxnwUq2XkHyG154vQ/viewform) first. See our [license file](./LICENSE) for more details.

# 🧑‍💻 We are hiring!

We've recently closed a [$38 million Series B funding round](https://techcrunch.com/2021/03/04/stream-raises-38m-as-its-chat-and-activity-feed-apis-power-communications-for-1b-users/) and we keep actively growing.
Our APIs are used by more than a billion end-users, and you'll have a chance to make a huge impact on the product within a team of the strongest engineers all over the world.

Check out our current openings and apply via [Stream's website](https://getstream.io/team/#jobs).
