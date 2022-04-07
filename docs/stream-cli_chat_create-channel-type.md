## stream-cli chat create-channel-type

Create channel type

### Synopsis

This command creates a new channel type. The 'properties' are raw JSON string.
The available properties can be found in the Go SDK's 'ChannelType' struct.


```
stream-cli chat create-channel-type --properties [raw-json-properties] [flags]
```

### Examples

```
# Create a new channel type called my-ch-type
$ stream-cli chat create-channel-type -p "{\"name\": \"my-ch-type\"}"

# Create a new channel type called reactionless with reactions disabled
$ stream-cli chat create-channel-type -p "{\"name\": \"reactionless\", \"reactions\": false}"

```

### Options

```
  -h, --help                help for create-channel-type
  -p, --properties string   [required] Raw JSON properties
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

