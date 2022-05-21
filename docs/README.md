# 📚 Documentation <!-- omit in toc -->

Stream's Command Line Interface (CLI) makes it easy to create and manage your [Stream](https://getstream.io) apps directly from the terminal.

> Currently, only Chat is supported; however, the ability to manage Feeds will be coming soon.

The generated CLI documentation is available [here](./stream-cli.md) - you can learn about all of the available commands there.

- [🏗 Installation](#-installation)
  - [Download the binaries](#download-the-binaries)
  - [Homebrew](#homebrew)
  - [Compile yourself](#compile-yourself)
- [🚀 Getting Started](#-getting-started)
- [📃 Use cases and examples](#-use-cases-and-examples)
- [🚨 Warning](#-warning)
- [🔨 Syntax](#-syntax)
- [💬 Auto completion](#-auto-completion)
- [🗒 Issues](#-issues)
- [📝 Changelog](#-changelog)

# 🏗 Installation

The Stream CLI is written in Go and precompiled into a single binary. It doesn't have any prerequisites.

## Download the binaries
You can find the binaries in the [Release section](https://github.com/GetStream/stream-cli/releases) of the repository. We also wrote a short script to download them and put it to your $PATH.

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
$ brew install stream-cli
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

![Stream](./first_config.svg)

> Note: Your API key and secret can be found on the [Stream Dashboard](https://getstream.io/dashboard) and is specific to your application.

# 📃 Use cases and examples

A couple of example use cases can be found [here](./use_cases.md). We've also created a separate documentation [for the import feature](./imports.md).

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

# 🗒 Issues

If you're experiencing problems directly related to the CLI, please add an [issue on GitHub](https://github.com/getstream/stream-cli/issues).

For other issues, submit a [support ticket](https://getstream.io/support).

# 📝 Changelog

As with any project, things are always changing. If you're interested in seeing what's changed in the Stream CLI, the changelog for this project can be tracked in the [Release](https://github.com/GetStream/stream-cli/releases) page of the repository.
