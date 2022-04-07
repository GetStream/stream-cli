## stream-cli chat list-imports

List imports

```
stream-cli chat list-imports --offset [int] --limit [int] --output-format [json|tree] [flags]
```

### Examples

```
# List all imports as json (default)
$ stream-cli chat list-imports

# List all imports as browsable tree
$ stream-cli chat list-imports --output-format tree

```

### Options

```
  -h, --help                   help for list-imports
  -l, --limit int              [optional] The number of imports returned (default 10)
  -O, --offset int             [optional] The starting offset of imports returned
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

