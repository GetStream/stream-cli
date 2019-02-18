![Stream Cli](https://i.imgur.com/H8AScTq.png)

# Stream CLI

Stream's Command Line Interface (CLI) makes it easy to create and manage your [Stream](https://getstream.io) apps directly from the terminal. Currently, only Chat is supported; however, the ability to manage Feeds will be coming soon.

[![Version](https://img.shields.io/npm/v/getstream-cli.svg)](https://npmjs.org/package/getstream-cli)
[![Dependency Status](https://david-dm.org/getstream/stream-cli/status.svg)](https://david-dm.org/getstream/stream-cli)
[![devDependency Status](https://david-dm.org/getstream/stream-cli/dev-status.svg)](https://david-dm.org/getstream/stream-cli?type=dev)
[![License](https://img.shields.io/npm/l/getstream-cli.svg)](https://github.com/getstream/stream-cli/blob/master/package.json)

## 🗒 Issues

If you're experiencing problems directly related to the CLI, please add an [issue on GitHub](https://github.com/getstream/stream-cli/issues).

For other issues, submit a [support ticket](https://getstream.io/support).

## 📚 Changelog

As with any project, things are always changing. If you're interested in seeing what's changed in the Stream CLI, the changelog for this project can be found [here](https://github.com/getstream/stream/blob/master/CHANGELOG.md).

## 🏗 Installation

The Stream CLI is easy to install. You have the option to use [homebrew](https://brew.sh) (preferred) if you're on macOS, download a single binary for your OS of choice, or download via npm.

#### Homebrew

```sh-session
$ brew install stream
```

#### NPM

```sh-session
$ npm install getstream-cli
```

#### Binaries

-   [Mac OS X](https://github.com/GetStream/stream-cli/releases)
-   [Linux](https://github.com/GetStream/stream-cli/releases)
-   [Windows](https://github.com/GetStream/stream-cli/releases)

> Note: Binaries are generally updated less frequently than Homebrew.

## 🚀 Getting Started

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

## 🔨 Commands

-   [`$ stream autocomplete`](#-stream-autocomplete)
-   [`$ stream commands`](#-stream-commands)
-   [`$ stream help`](#-stream-help)
-   stream config
    -   [`$ set`](#-stream-configset)
    -   [`$ get`](#-stream-configget)
    -   [`$ destroy`](#-stream-configdestroy)
-   stream channel
    -   [`$ edit`](#-stream-channeledit)
    -   [`$ get`](#-stream-channelget)
    -   [`$ init`](#-stream-channelinit)
    -   [`$ list`](#-stream-channellist)
    -   [`$ query`](#-stream-channelquery)
-   stream message
    -   [`$ send`](#-stream-messagesend)
    -   [`$ remove`](#-stream-messageremove)
-   stream moderate
    -   [`$ ban`](#-stream-moderateban)
    -   [`$ flag`](#-stream-moderateflag)
    -   [`$ mute`](#-stream-moderatemute)
-   stream user
    -   [`$ add`](#-stream-useradd)
    -   [`$ remove`](#-stream-userremove)

### `$ stream autocomplete`

Initialize autocomplete for the CLI **(recommended)**

```sh-session
$ stream autocomplete
```

### `$ stream commands`

Display all commands

```sh-session
$ stream commands
```

### `$ stream help`

Get help with the CLI

```sh-session
$ stream help
```

### `$ stream config:set`

Initialize a new configuration file.

```sh-session
USAGE
  $ stream config:set
```

### `$ stream config:get`

Retrieve your configuration settings.

```sh-session
USAGE
  $ stream config:get
```

### `$ stream config:destroy`

Destroy your configuration file

```sh-session
USAGE
  $ stream config:destroy
```

> Note: The command `stream config:set` must be called to re-initialize the configuration.

### `$ stream channel:edit`

Edit a specified channel

```sh-session
USAGE
  $ stream channel:edit

OPTIONS
  -d, --data=data                                       Additional data as a JSON payload.
  -i, --id=id                                           (required) Channel ID.
  -m, --members=members                                 Comma separated list of members.
  -n, --name=name                                       (required) Name of room.
  -r, --reason=reason                                   (required) Reason for changing channel.
  -t, --type=livestream|messaging|gaming|commerce|team  (required) Type of channel.
  -u, --url=url                                         URL to channel image.
```

### `$ stream channel:get`

Get a specified channel

```sh-session
USAGE
  $ stream channel:get

OPTIONS
  -i, --id=id                                           (required) Channel ID.
  -t, --type=livestream|messaging|gaming|commerce|team  (required) Type of channel.
```

### `$ stream channel:init`

Initialize a new channel

```sh-session
USAGE
  $ stream channel:init

OPTIONS
  -d, --data=data                                       Additional data as a JSON payload.
  -i, --id=id                                           (required) [default: <UUID>] Channel ID.
  -m, --members=members                                 Comma separated list of members.
  -n, --name=name                                       (required) Name of room.
  -t, --type=livestream|messaging|gaming|commerce|team  (required) Type of channel.
  -u, --image=image                                     URL to channel image.
```

### `$ stream channel:list`

List all channels associated with your account

```sh-session
USAGE
  $ stream channel:list
```

### `$ stream channel:query`

Query for channels

```sh-session
USAGE
  $ stream channel:query

OPTIONS
  -f, --filter=filter                                   Filters to apply.
  -i, --id=id                                           [default: <UUID>] Channel ID.
  -s, --sort=sort                                       Sort to apply.
  -t, --type=livestream|messaging|gaming|commerce|team  Type of channel.
```

### `$ stream message:send`

Send a message to a specific channel

```sh-session
USAGE
  $ stream message:send

OPTIONS
  -a, --attachments=attachments                         JSON payload of attachments
  -i, --id=id                                           [default: <UUID>] Channel ID.
  -m, --message=message                                 (required) Message to send.
  -t, --type=livestream|messaging|gaming|commerce|team  (required) Type of channel.
  -u, --user=user                                       (required) [default: *] ID of user.
```

### `$ stream message:remove`

Remove a message from a channel

```sh-session
USAGE
  $ stream message:remove

OPTIONS
  -i, --id=id  (required) Channel ID.
```

### `$ stream moderate:ban`

Ban a user from a channel forever or based on a per minute timeout

```sh-session
USAGE
  $ stream moderate:ban

OPTIONS
  -r, --reason=reason    (required) Reason for timeout.
  -t, --timeout=timeout  (required) [default: 60] Timeout in minutes.
  -u, --user=user        (required) ID of user.
```

### `$ stream moderate:flag`

Flag users and messages for inappropriate behavior or explicit content

```sh-session
USAGE
  $ stream moderate:flag

OPTIONS
  -m, --message=message  ID of message.
  -u, --user=user        ID of user.
```

### `$ stream moderate:mute`

Mute a user in a channel

```sh-session
USAGE
  $ stream moderate:mute

OPTIONS
  -u, --user=user  (required) User ID.
```

### `$ stream user:add`

Add a user to a channel and specify permissions

```sh-session
USAGE
  $ stream user:add

OPTIONS
  -i, --id=id                  (required) Channel name.
  -m, --moderators=moderators  (required) Comma separated list of moderators to add.
  -t, --type=type              (required) Channel type.
```

### `$ stream user:remove`

Remove a user from a channel

```sh-session
USAGE
  $ stream user:remove

OPTIONS
  -i, --id=id                  (required) Channel name.
  -m, --moderators=moderators  (required) Comma separated list of moderators to remove.
  -t, --type=type              (required) Channel type.
```

## 📣 Feedback

If you have any suggestions or just want to let us know what you think of the Stream CLI, please send us a message at support@getstream.io or create a [GitHub Issue](https://github.com/getstream/stream-cli/issues).
