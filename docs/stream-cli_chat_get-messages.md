## stream-cli chat get-messages

Return multiple messages

```
stream-cli chat get-messages --channel-type [channel-type] --channel-id [channel-id] --output-format [json|tree] [message-id-1] [message-id-2] [message-id ...] [flags]
```

### Examples

```
# Returns 3 messages of 'redteam' channel of 'messaging' channel type
$ stream-cli chat get-messages --channel-type messaging --channel-id redteam msgid-1 msgid-2 msgid-3

```

### Options

```
  -i, --channel-id string      [required] Channel id
  -t, --channel-type string    [required] Channel type such as 'messaging' or 'livestream'
  -h, --help                   help for get-messages
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

