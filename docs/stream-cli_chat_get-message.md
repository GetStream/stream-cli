## stream-cli chat get-message

Return a single message

```
stream-cli chat get-message [message-id] --output-format [json|tree] [flags]
```

### Examples

```
# Returns a message with id 'msgid-1'
$ stream-cli chat get-message msgid-1

```

### Options

```
  -h, --help                   help for get-message
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

