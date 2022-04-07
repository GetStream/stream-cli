## stream-cli chat update-channel-partial

Update a channel partially

### Synopsis

Updates an existing channel. The 'set' property is a comma separated list of key value pairs.
The 'unset' property is a comma separated list of property names.


```
stream-cli chat update-channel-partial --type [channel-type] --id [channel-id] --set [key-value-pairs] --unset [property-names] [flags]
```

### Examples

```
# Freeze a channel and set 'age' to 21. At the same time, remove 'haircolor' and 'height'.
stream-cli chat update-channel-partial --type messaging --id channel1 --set frozen=true,age=21 --unset haircolor,height

```

### Options

```
  -h, --help                 help for update-channel-partial
  -i, --id string            [required] Channel id
  -s, --set stringToString   [optional] Comma-separated key-value pairs to set (default [])
  -t, --type string          [required] Channel type such as 'messaging' or 'livestream'
  -u, --unset string         [optional] Comma separated list of properties to unset
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

