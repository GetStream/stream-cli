## stream-cli import feeds get-import

Get a feeds import task

```
stream-cli import feeds get-import [task-id] --output-format [json|tree] --watch [flags]
```

### Examples

```
# Returns a feeds import and prints it as JSON
$ stream-cli import feeds get-import dcb6e366-93ec-4e52-af6f-b0c030ad5272

# Returns a feeds import and watches for completion
$ stream-cli import feeds get-import dcb6e366-93ec-4e52-af6f-b0c030ad5272 --watch

```

### Options

```
  -h, --help                   help for get-import
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
  -w, --watch                  [optional] Keep polling the import to track its status
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli import feeds](stream-cli_import_feeds.md)	 - Import data into Feeds

