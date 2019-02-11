# Stream Chat CLI

Stream's Chat Command Line Interface (CLI) makes it easy to create and manage your Stream Chat apps directly from the terminal.

[![Version](https://img.shields.io/npm/v/stream-chat-cli.svg)](https://npmjs.org/package/stream-chat-cli)
[![Downloads/week](https://img.shields.io/npm/dw/stream-chat-cli.svg)](https://npmjs.org/package/stream-chat-cli)
[![License](https://img.shields.io/npm/l/stream-chat-cli.svg)](https://github.com/nparsons08/stream-chat-cli/blob/master/package.json)

# Links

-   [Usage](#usage)
-   [Commands](#commands)

# Usage

```sh-session
$ npm install -g stream-chat-cli
$ chat COMMAND
running command...
$ chat (-v|--version|version)
stream-chat-cli/0.0.1 darwin-x64 node-v11.9.0
$ chat --help [COMMAND]
USAGE
  $ chat COMMAND
...
```

# Commands

-   [`chat init`](#stream-chat-init)
-   [`chat message [COMMAND]`](#stream-chat-cli-message)
-   [`chat watch [COMMAND]`](#stream-chat-cli-watch)

## `chat init`

The `chat init` command stores credentials for the the Stream Chat client. This allows later calls to authenticate against the credentials stored in the credentials configuration. The `chat init` command should be called prior to running any other command.

```
USAGE
  $ chat init

OPTIONS
  -v, --verify=true see JSON output of your stored credentials
  -d, --destroy=true destroy existing credentials

DESCRIPTION
  ...
  Initializes the CLI with API credentials for calls to the Stream API.
```

_See code: [src/commands/init.js](https://github.com/nparsons08/stream-chat/blob/v0.0.1/src/commands/hello.js)_

## `chat help [COMMAND]`

display help for stream-chat

```
USAGE
  $ chat help [COMMAND]

ARGUMENTS
  COMMAND  command to show help for

OPTIONS
  --all see all commands in CLI
```
