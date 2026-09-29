## stream-cli import feeds list-imports

List feeds import tasks

```
stream-cli import feeds list-imports --output-format [json|tree] --state [1-4] [flags]
```

### Examples

```
# List all feeds imports as json (default)
$ stream-cli import feeds list-imports

# List feeds imports filtered by state
$ stream-cli import feeds list-imports --state 2

# List all feeds imports as browsable tree
$ stream-cli import feeds list-imports --output-format tree

```

### Options

```
  -h, --help                   help for list-imports
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
  -s, --state int              [optional] Filter imports by state (1-4)
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli import feeds](stream-cli_import_feeds.md)	 - Import data into Feeds

