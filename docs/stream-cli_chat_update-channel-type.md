## stream-cli chat update-channel-type

Update channel type

### Synopsis

This command updates an existing channel type. The 'properties' are raw JSON string.
The available fields can be checked here:
https://getstream.io/chat/docs/rest/#channel-types-updatechanneltype


```
stream-cli chat update-channel-type --type [channel-type] --properties [raw-json-properties] [flags]
```

### Examples

```
# Enabling quotes in an existing channel type
$ stream-cli chat update-channel-type --type my-channel-type --properties '{"quotes": true}'

```

### Options

```
  -h, --help                help for update-channel-type
  -p, --properties string   [required] Raw JSON properties
  -t, --type string         [required] Channel type
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

