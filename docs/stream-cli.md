## stream-cli

Stream CLI

### Synopsis

Interact with your Stream applications easily

### Examples

```
# Get Chat application settings
$ stream-cli chat get-app

# List all Chat channel types
$ stream-cli chat list-channel-types

# Create a new Chat user
$ stream-cli chat upsert-user --properties "{\"id\":\"my-user-1\"}"

```

### Options

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
  -h, --help            help for stream-cli
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications
* [stream-cli config](stream-cli_config.md)	 - Manage app configurations

