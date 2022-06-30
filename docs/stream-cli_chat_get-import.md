## stream-cli chat get-import

Get import

```
stream-cli chat get-import [import-id] --output-format [json|tree] --watch [flags]
```

### Examples

```
# Returns an import and prints it as JSON
$ stream-cli chat get-import dcb6e366-93ec-4e52-af6f-b0c030ad5272

# Returns an import and prints it as JSON, and wait for it to complete
$ stream-cli chat get-import dcb6e366-93ec-4e52-af6f-b0c030ad5272 --watch

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

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

