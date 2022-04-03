## stream-cli chat get-channel

Return a channel

```
stream-cli chat get-channel --type [channel-type] --id [channel-id] [flags]
```

### Examples

```
# Returns 'redteam' channel of 'messaging' channel type as JSON
$ stream-cli chat get-channel --type messaging --id redteam

# Returns 'blueteam' channel of 'messaging' channel type as a browsable tree
$ stream-cli chat get-channel --type messaging --id blueteam --output-format tree

```

### Options

```
  -h, --help                   help for get-channel
  -i, --id string              [required] Channel id
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
  -t, --type string            [required] Channel type such as 'messaging' or 'livestream'
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

