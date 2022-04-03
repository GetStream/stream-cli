## stream-cli chat get-channel-type

Get channel type

```
stream-cli chat get-channel-type [channel-type] --output-format [json|tree] [flags]
```

### Examples

```
# Returns a channel type and prints it as JSON
$ stream-cli chat get-channel-type livestream

# Returns a channel type and prints it as a browsable tree
$ stream-cli chat get-channel-type messaging --output-format tree

```

### Options

```
  -h, --help                   help for get-channel-type
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

