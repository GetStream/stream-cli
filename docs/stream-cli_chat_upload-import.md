## stream-cli chat upload-import

Upload an import

```
stream-cli chat upload-import [filename] --mode [upsert|insert] --output-format [json|tree] [flags]
```

### Examples

```
# Uploads an import and prints it as JSON
$ stream-cli chat upload-import data.json --mode insert

# Uploads an import and prints it as a browsable tree
$ stream-cli chat upload-import data.json --mode insert --output-format tree

```

### Options

```
  -h, --help                   help for upload-import
  -m, --mode string            [optional] Import mode. Canbe upsert or insert (default "upsert")
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

