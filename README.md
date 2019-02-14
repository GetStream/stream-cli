![Stream Cli](https://i.imgur.com/H8AScTq.png)

# Stream Chat CLI

Stream's Command Line Interface (CLI) makes it easy to create and manage your Stream apps directly from the terminal. Currently, only Chat is supported; however, the ability to manage Feeds will be coming soon.

[![Version](https://img.shields.io/npm/v/stream-cli.svg)](https://npmjs.org/package/getstream-cli)
[![Dependency Status](https://david-dm.org/getstream/stream-cli/status.svg)](https://david-dm.org/getstream/stream-cli)
[![devDependency Status](https://david-dm.org/getstream/stream-cli/dev-status.svg)](https://david-dm.org/getstream/stream-cli?type=dev)
[![License](https://img.shields.io/npm/l/stream-cli.svg)](https://github.com/getstream/stream-cli/blob/master/package.json)

# 📌 Requirements

Only Node 8+ is supported. Node 6 will reach end-of-life April 2019. At that point we will continue to support the current LTS version of node. You can add the node package to your CLI to ensure users are on Node 8.

# 🗒 Issues

If you're experiencing problems directly related to the CLI, please add an [issue on GitHub](https://github.com/getstream/stream-cli/issues).

For other issues, submit a [support ticket](https://getstream.io/support).

# 🏗 Installation

The Stream CLI is easy to install. You have the option to download a single binary (preferred) with zero run-time dependencies for your OS of choice, or install it using [NPM](https://www.npmjs.com/package/getstream-cli).

### Binaries

-   [Mac OS X](https://github.com/GetStream/stream-cli/releases)
-   [Linux](https://github.com/GetStream/stream-cli/releases)
-   [Windows](https://github.com/GetStream/stream-cli/releases)

### NPM

```sh-session
$ npm install -g getstream-cli
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

# 🔨 Commands

-   [stream autocomplete](#$-stream-autocomplete)
-   [stream commands](#)
-   [stream help](#)
-   [stream config](#)
    -   [set](#)
    -   [get](#)
    -   [destroy](#)
-   [stream channel](#)
    -   [edit](#)
    -   [get](#)
    -   [init](#)
    -   [list](#)
    -   [query](#)
-   [stream message](#)
    -   [send](#)
    -   [remove](#)
-   [stream moderate](#)
    -   [ban](#)
    -   [flag](#)
    -   [mute](#)
-   [stream user](#)
    -   [add](#)
    -   [ban](#)
    -   [remove](#)

## `$ stream autocomplete`

### Input

```sh-session
$ stream autocomplete
```

### Output

```sh-session
Building the autocomplete cache... done

Setup Instructions for STREAM CLI Autocomplete ---

1) Add the autocomplete env var to your zsh profile and source it
$ printf "$(stream autocomplete:script zsh)" >> ~/.zshrc; source ~/.zshrc

NOTE: After sourcing, you can run `$ compaudit -D` to ensure no permissions conflicts are present

2) Test it out, e.g.:
$ stream <TAB>                 # Command completion
$ stream command --<TAB>       # Flag completion

Enjoy!
```

## `$ stream commands`

### Input

```sh-session
$ stream commands
```

### Output

```sh-session
autocomplete
channel:edit
channel:get
channel:init
channel:list
channel:query
commands
config:destroy
config:get
config:set
help
log
message:remove
message:send
moderate:ban
moderate:flag
moderate:mute
user:add
user:ban
user:remove
```
