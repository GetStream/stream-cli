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
You can find the binaries in the [Release section](https://github.com/GetStream/stream-cli/releases) of this repository.

<details><summary>One liners for downloading the executable</summary>

<details markdown="1"><summary><strong>MacOS</strong></summary>

## **ARM** <!-- omit in toc -->
```shell
$ export URL=$(curl -s https://api.github.com/repos/GetStream/stream-cli/releases/latest | grep Darwin_arm  | cut -d '"' -f 4 | sed '1d')
$ curl -L $URL -o stream-cli.tar.gz
$ tar -xvf stream-cli.tar.gz
```

## **Intel** <!-- omit in toc -->
```shell
$ export URL=$(curl -s https://api.github.com/repos/GetStream/stream-cli/releases/latest | grep Darwin_x86  | cut -d '"' -f 4 | sed '1d')
$ curl -L $URL -o stream-cli.tar.gz
$ tar -xvf stream-cli.tar.gz
```

You can either put it to your $PATH or set up a symbolic link:
```shell
$ ln -s $PWD/stream-cli /usr/local/bin/stream-cli
```

</details>

<details markdown="1"><summary><strong>Linux</strong></summary>

## **ARM** <!-- omit in toc -->
```shell
$ export URL=$(curl -s https://api.github.com/repos/GetStream/stream-cli/releases/latest | grep Linux_arm64  | cut -d '"' -f 4 | sed '1d')
$ curl -L $URL -o stream-cli.tar.gz
$ tar -xvf stream-cli.tar.gz
```

## **Intel** <!-- omit in toc -->
```shell
$ export URL=$(curl -s https://api.github.com/repos/GetStream/stream-cli/releases/latest | grep Linux_x86  | cut -d '"' -f 4 | sed '1d')
$ curl -L $URL -o stream-cli.tar.gz
$ tar -xvf stream-cli.tar.gz
```

You can either put it to your $PATH or set up a symbolic link:
```shell
$ ln -s $PWD/stream-cli /usr/local/bin/stream-cli
```

</details>
<details markdown="1"><summary><strong>Windows</strong></summary>

## **ARM** <!-- omit in toc -->
```powershell
> $latestRelease = Invoke-WebRequest "https://api.github.com/repos/GetStream/stream-cli/releases/latest"
> $json = $latestRelease.Content | ConvertFrom-Json
> $url = $json.assets | ? { $_.name -match "Windows_arm" } | select -expand browser_download_url
> Invoke-WebRequest -Uri $url -OutFile "stream-cli.zip"
> Expand-Archive -Path ".\stream-cli.zip"
```

## **Intel** <!-- omit in toc -->
```powershell
> $latestRelease = Invoke-WebRequest "https://api.github.com/repos/GetStream/stream-cli/releases/latest"
> $json = $latestRelease.Content | ConvertFrom-Json
> $url = $json.assets | ? { $_.name -match "Windows_x86" } | select -expand browser_download_url
> Invoke-WebRequest -Uri $url -OutFile "stream-cli.zip"
> Expand-Archive -Path ".\stream-cli.zip"
```

</details>
</details>


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

A couple of example use cases can be found [here](./use_cases.md).

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
