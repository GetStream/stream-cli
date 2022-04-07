## stream-cli chat update-channel

Update a channel

### Synopsis

Updates an existing channel. The 'properties' are specified as a raw json string. The valid
properties are the 'ChannelRequest' object of the official documentation.
Such as 'team', 'frozen', 'disabled' or any custom property.
https://getstream.io/chat/docs/rest/#channels-updatechannel


```
stream-cli chat update-channel --type [channel-type] --id [channel-id] --properties [raw-json-properties] [flags]
```

### Examples

```
# Unfreeze a channel
$ stream-cli chat update-channel --type messaging --id redteam --properties "{\"frozen\":false}"

```

### Options

```
  -h, --help                help for update-channel
  -i, --id string           [required] Channel id
  -p, --properties string   [required] Channel properties to update
  -t, --type string         [required] Channel type such as 'messaging' or 'livestream'
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

