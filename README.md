![Stream Cli](https://i.imgur.com/H8AScTq.png)

# Stream CLI

> Note: The Stream CLI is currently in beta and may contain bugs. This _should not_ be used against a production environment at this time. To report bugs, please follow the instructions below. Thank you for your support!

Stream's Command Line Interface (CLI) makes it easy to create and manage your [Stream](https://getstream.io) apps directly from the terminal. Currently, only Chat is supported; however, the ability to manage Feeds will be coming soon.

[![Version](https://img.shields.io/npm/v/getstream-cli.svg)](https://npmjs.org/package/getstream-cli)
[![Dependency Status](https://david-dm.org/getstream/stream-cli/status.svg)](https://david-dm.org/getstream/stream-cli)
[![devDependency Status](https://david-dm.org/getstream/stream-cli/dev-status.svg)](https://david-dm.org/getstream/stream-cli?type=dev)
[![License](https://img.shields.io/npm/l/getstream-cli.svg)](https://github.com/getstream/stream-cli/blob/master/package.json)

# 🗒 Issues

If you're experiencing problems directly related to the CLI, please add an [issue on GitHub](https://github.com/getstream/stream-cli/issues).

For other issues, submit a [support ticket](https://getstream.io/support).

# 📚 Changelog

As with any project, things are always changing. If you're interested in seeing what's changed in the Stream CLI, the changelog for this project can be found [here](https://github.com/getstream/stream/blob/master/CHANGELOG.md).

# 🏗 Installation

The Stream CLI is easy to install and available via npm.

```sh-session
$ npm install -g getstream-cli
```

**OR**

```sh-session
$ yarn global add getstream-cli
```

# 🚀 Getting Started

In order to initialize the CLI, please have your Stream API key and secret ready. Run the following command:

```sh-session
$ stream config:set
```

You will then be prompted to enter your API key and secret.

```sh-session
$ ? What's your API key? 🔒
$ ? What's your API secret? 🔒
```

Now, you're good to go!

```sh-session
$ Your config has been generated! 🚀
```

> Note: Your API key and secret can be found on the [Stream Dashboard](https://getstream.io/dashboard) and is specific to your application.

# 🔨 Syntax

Basic commands use the following syntax:

```sh-session
$ stream command:COMMAND --arg1 "foo" --arg2 "bar"
```

Whereas commands for specific products use subcommands:

```sh-session
$ stream command:COMMAND:SUBCOMMAND --arg1 "foo" --arg2 "bar"
```

# 🤔 Fun Facts

Interested in using the calling the CLI from a script? You can do that! Many of the commands accept a `raw` argument as a `boolean`. Just pass the following along to the CLI and you'll be able to get a full representation of the response data behind call:

```sh-session
$ stream command:COMMAND --arg1 "foo" --arg2 "bar" --raw "true"
```

Need to copy the output to your clipboard? Not a problem. Just pipe the information to `pbcopy` (on macOS) along with the `--raw` argument set to `true`:

```sh-session
$ stream debug:token --jwt "foo.bar.baz" --raw "true" | pbcopy
```

# 🥳‍ Usage

<!-- usage -->

```sh-session
$ npm install -g getstream-cli
$ stream COMMAND
running command...
$ stream (-v|--version|version)
getstream-cli/0.0.1-beta.22 darwin-x64 node-v10.15.1
$ stream --help [COMMAND]
USAGE
  $ stream COMMAND
...
```

<!-- usagestop -->

# 💻 Commands

<!-- commands -->

-   [`stream autocomplete`](docs/autocomplete.md) - display autocomplete installation instructions
-   [`stream chat`](docs/chat.md) - configure and manage all things related to chat
-   [`stream commands`](docs/commands.md) - list all the commands
-   [`stream config`](docs/config.md) - manage configuration variables
-   [`stream debug`](docs/debug.md) - debugging tools for the command line
-   [`stream help`](docs/help.md) - display help for stream

<!-- commandsstop -->

# 📣 Feedback

If you have any suggestions or just want to let us know what you think of the Stream CLI, please send us a message at support@getstream.io or create a [GitHub Issue](https://github.com/getstream/stream-cli/issues).
