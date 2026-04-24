## stream-cli import feeds upload-import

Upload an import for Feeds

```
stream-cli import feeds upload-import [filename] --output-format [json|tree] [flags]
```

### Examples

```
# Uploads a feeds import and prints it as JSON
$ stream-cli import feeds upload-import data.json

# Uploads a feeds import and prints it as a browsable tree
$ stream-cli import feeds upload-import data.json --output-format tree

```

### Options

```
  -h, --help                    help for upload-import
  -o, --output-format string    [optional] Output format. Can be json or tree (default "json")
      --skip-references-check   [optional] Skip references validation for the import (default false)
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli import feeds](stream-cli_import_feeds.md)	 - Import data into Feeds

