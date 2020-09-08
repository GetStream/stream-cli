![Stream Cli](https://i.imgur.com/H8AScTq.png)

# Stream CLI

Stream's Command Line Interface (CLI) makes it easy to create and manage your [Stream](https://getstream.io) apps directly from the terminal. Currently, only Chat is supported; however, the ability to manage Feeds will be coming soon.

> **Note**: The latest version of Node (v14.x) is required for this CLI. If you are looking for a way to manage multiple Node.js environments, please look at [nvm](https://github.com/nvm-sh/nvm).

[![Coverage Status](https://coveralls.io/repos/github/GetStream/stream-cli/badge.svg?branch=master)](https://coveralls.io/github/GetStream/stream-cli?branch=master)
[![Version](https://img.shields.io/npm/v/getstream-cli.svg)](https://npmjs.org/package/getstream-cli)
[![License](https://img.shields.io/npm/l/getstream-cli.svg)](https://github.com/getstream/stream-cli/blob/master/package.json)

# 🗒 Issues

If you're experiencing problems directly related to the CLI, please add an [issue on GitHub](https://github.com/getstream/stream-cli/issues).

For other issues, submit a [support ticket](https://getstream.io/support).

# 📚 Changelog

As with any project, things are always changing. If you're interested in seeing what's changed in the Stream CLI, the changelog for this project can be found [here](https://github.com/getstream/stream/blob/master/CHANGELOG.md).

# 🏗 Installation

The Stream CLI is easy to install and available via npm. The CLI requires Node v10.x or above.

```sh-session
$ yarn global add getstream-cli
```

**OR**

```sh-session
$ npm install -g getstream-cli
```

# 🚀 Getting Started

In order to initialize the CLI, it's as simple as:

![Stream](https://i.imgur.com/SA9uMQ1.png)

> Note: Your API key and secret can be found on the [Stream Dashboard](https://getstream.io/dashboard) and is specific to your organization.

# 🔨 Syntax

Basic commands use the following syntax:

```sh-session
$ stream command:COMMAND --arg1 "foo" --arg2 "bar"
```

Whereas commands for specific products use subcommands:

```sh-session
$ stream command:COMMAND:SUBCOMMAND --arg1 "foo" --arg2 "bar"
```

# 🎩 Fun Facts

Interested in using the calling the CLI from a script? Or maybe you simply want raw response data? You can do that! Many of the commands accept a `json` argument as a `boolean`. Just pass the following along to the CLI and you'll get back a full representation of the response (in a raw data format):

```sh-session
$ stream command:COMMAND --arg1 "foo" --arg2 "bar" --json
```

Need to copy the output to your clipboard? Not a problem. Just pipe the information to `pbcopy` (on macOS) along with the `--json` flag:

```sh-session
$ stream debug:token --token "foo.bar.baz" --json | pbcopy
```

Want to call the Stream CLI using a bash command? No problem! Make sense of output using [jq](https://stedolan.github.io/jq/), a lightweight and flexible command-line JSON processor.

```bash
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
```

**OR**

```bash
#! /bin/bash

stream chat:channel:create --channel=$(openssl rand -hex 12) --type="messaging" --name="CLI" --json | jq '.'
```

> Note: See [here](https://github.com/GetStream/stream-cli/tree/master/examples/bash) for additional examples!

# 🥳‍ Usage

<!-- usage -->
```sh-session
$ npm install -g getstream-cli
$ stream COMMAND
running command...
$ stream (-v|--version|version)
getstream-cli/0.0.62 darwin-x64 node-v14.8.0
$ stream --help [COMMAND]
USAGE
  $ stream COMMAND
...
```
<!-- usagestop -->

# 💻 Commands

<!-- commands -->
# Command Topics

* [`stream autocomplete`](docs/autocomplete.md) - display autocomplete installation instructions
* [`stream chat`](docs/chat.md) - Adds a member to a channel.
* [`stream commands`](docs/commands.md) - list all the commands
* [`stream config`](docs/config.md) - Destroys your user configuration.
* [`stream debug`](docs/debug.md) - Debugs a JWT token provided by Stream.
* [`stream help`](docs/help.md) - display help for stream

<!-- commandsstop -->

# 📣 Feedback

If you have any suggestions or just want to let us know what you think of the Stream CLI, please send us a message at support@getstream.io or create a [GitHub Issue](https://github.com/getstream/stream-cli/issues).
