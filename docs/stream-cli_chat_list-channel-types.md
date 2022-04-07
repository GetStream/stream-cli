## stream-cli chat list-channel-types

List channel types

### Synopsis

This command lists all channel types, including built-in and custom ones.


```
stream-cli chat list-channel-types --output-format [json|tree] [flags]
```

### Examples

```
# List all channel types as json (default)
$ stream-cli chat list-channel-types

# List all channel types as browsable tree
$ stream-cli chat list-channel-types --output-format tree

```

### Options

```
  -h, --help                   help for list-channel-types
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

