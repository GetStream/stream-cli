![Stream Cli](https://i.imgur.com/H8AScTq.png)

# Stream CLI

Stream's Command Line Interface (CLI) makes it easy to create and manage your [Stream](https://getstream.io) apps directly from the terminal. Currently, only Chat is supported; however, the ability to manage Feeds will be coming soon.

> **Note**: The latest version of Node (v14.x) is required for this CLI. If you are looking for a way to manage multiple Node.js environments, please look at [nvm](https://github.com/nvm-sh/nvm).



# 🏗 Installation

The Stream CLI is built with Go and we compile it into a single binary so you don't need any prerequisite available on your computer. You can simply download it and put it to your `$PATH`.

The releases are available in [GitHub Releases](https://github.com/GetStream/stream-cli/releases).

On MacOS:
```shell
$ curl -L https://github.com/GetStream/stream-cli/releases/latest/download/darwin -o $HOME/stream-cli
$ sudo ln -s $HOME/stream-cli /usr/local/bin/stream-cli

# Check if it works:
$ stream-cli --version
```

On Linux:
```shell
$ curl -L https://github.com/GetStream/stream-cli/releases/latest/download/linux -o $HOME/stream-cli
$ sudo ln -s $HOME/stream-cli /usr/local/bin/stream-cli

# Check if it works:
$ stream-cli --version
```

> 💡 Note for Linux and MacOS users
>
> Unfortunately, the `stream` binary name is already [used by an app called ImageMagick](https://github.com/GetStream/stream-cli/issues/63) hence we opted to use the name `stream-cli`. If you are not using ImageMagick, you can just use `stream` name instead of `stream-cli`:
> ```shell
> $ ln -s $HOME/stream-cli /usr/local/bin/stream
> $ stream --version
> ```

On Windows:
```powershell
$ Invoke-WebRequest "https://github.com/GetStream/stream-cli/releases/latest/download/windows.exe" -Out "$env:APPDATA/stream-cli/stream-cli.exe"
$ [Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:APPDATA/stream-cli", "User")

# Check if it works
$ stream-cli --version

# Set an alias in PowerShell if you want to
$ Set-Alias -Name stream -Value "$env:APPDATA/stream-cli/stream-cli.exe" -Scope "Global"
$ stream --version
```

# 🪄 Self update

The CLI has a self-update mechanism where it automatically downloads the latest version from GitHub Releases and overwrites the current binary.

```shell
$ stream update self

> Successfully updated to v1.5.0 ✅
```

# 🚀 Getting Started

In order to initialize the CLI, it's as simple as:

![Init](./assets/stream_init.svg)

> Note: Your API key and secret can be found on the [Stream Dashboard](https://getstream.io/dashboard) and is specific to your organization.

# 🔨 Syntax

Basic commands use the following syntax:

```shell
$ stream command --arg1 "foo" --arg2 "bar"
```

Whereas commands for specific products use subcommands:

```shell
$ stream command subcommand --arg1 "foo" --arg2 "bar"
```

# 🎩 Fun Facts

Interested in using the calling the CLI from a script? Or maybe you simply want raw response data? You can do that! Many of the commands accept a `json` argument as a `boolean`. Just pass the following along to the CLI and you'll get back a full representation of the response (in a raw data format):

```shell
$ stream command:COMMAND --arg1 "foo" --arg2 "bar" --json
```

Need to copy the output to your clipboard? Not a problem. Just pipe the information to `pbcopy` (on macOS) along with the `--json` flag:

```shell
$ stream channel list -t messaging --json| pbcopy
```

Want to call the Stream CLI using a bash command? No problem! Make sense of output using [jq](https://stedolan.github.io/jq/), a lightweight and flexible command-line JSON processor.

```bash
#! /bin/bash

run=$(stream config get --json)

name=$(jq --raw-output '.name' <<< "${run}")
email=$(jq --raw-output '.email' <<< "${run}")
apiKey=$(jq --raw-output '.apiKey' <<< "${run}")
apiSecret=$(jq --raw-output '.apiSecret' <<< "${run}")

echo $name
echo $email
echo $apiKey
echo $apiSecret
```

**OR**

```bash
#! /bin/bash

stream chat:channel:create --channel=$(openssl rand -hex 12) --type="messaging" --name="CLI" --json | jq '.'
```

> Note: See [here](https://github.com/GetStream/stream-cli/tree/master/examples/bash) for additional examples!

# 💻 Commands

<!-- commands -->
## Command Topics

* [`stream autocomplete`](docs/autocomplete.md) - display autocomplete installation instructions
* [`stream chat`](docs/chat.md) - Manage chat
* [`stream commands`](docs/commands.md) - list all the commands
* [`stream config`](docs/config.md) - Configure API access
* [`stream debug`](docs/debug.md) - Debugging tools
* [`stream help`](docs/help.md) - display help for stream

<!-- commandsstop -->

# 📣 Feedback

If you have any suggestions or just want to let us know what you think of the Stream CLI, please send us a message at support@getstream.io or create a [GitHub Issue](https://github.com/getstream/stream-cli/issues).

# 🗒 Issues

If you're experiencing problems directly related to the CLI, please add an [issue on GitHub](https://github.com/getstream/stream-cli/issues).

For other issues, submit a [support ticket](https://getstream.io/support).

# 📚 Changelog

As with any project, things are always changing. If you're interested in seeing what's changed in the Stream CLI, the changelog for this project can be found [here](https://github.com/getstream/stream-cli/blob/master/CHANGELOG.md).
# 🔧 Development

This project contains generated code and documentation. In order to apply changes you should run the following command:

```bash
$ yarn run generate
```