## stream-cli chat list-channels

List channels

### Synopsis

List all channels of a given channel type. You can also provide
a limit for paginating the results.


```
stream-cli chat list-channels --type [channel-type] [flags]
```

### Examples

```
# List the top 5 'messaging' channels as a json
$ stream-cli chat list-channels --type messaging --limit 5

# List the top 20 'livestream' channels as a browsable tree
$ stream-cli chat list-channels --type livestream --limit 20 --output-format tree

```

### Options

```
  -h, --help                   help for list-channels
  -l, --limit int              [optional] Number of channels to return. Used for pagination (default 10)
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

