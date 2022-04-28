## stream-cli config remove

Remove one or more application.

### Synopsis

Remove one or more application from the configuraiton file. This operation is irrevocable.

```
stream-cli config remove [app-name-1] [app-name-2] [app-name-n] [flags]
```

### Examples

```
# Remove a single application from the CLI
$ stream-cli config remove staging

# Remove multiple applications from the CLI
$ stream-cli config remove staging testing

```

### Options

```
  -h, --help   help for remove
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli config](stream-cli_config.md)	 - Manage app configurations

