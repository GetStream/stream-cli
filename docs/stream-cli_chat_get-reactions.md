## stream-cli chat get-reactions

Get reactions for a message

```
stream-cli chat get-reactions [message-id] [flags]
```

### Examples

```
# Get reactions for a [08f64828-3bba-42bd-8430-c26a3634ee5c] message
$ stream-cli chat get-reactions 08f64828-3bba-42bd-8430-c26a3634ee5c --output-format json

```

### Options

```
  -h, --help                   help for get-reactions
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

