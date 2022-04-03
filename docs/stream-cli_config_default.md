## stream-cli config default

Set an application as the default

### Synopsis

Set an application as the default which will be used
for all further operations unless specified otherwise.


```
stream-cli config default [app-name] [flags]
```

### Examples

```
# Set an application as the default
$ stream-cli config default staging

# All underlying operations will use it if not specified otherwise
$ stream-cli chat get-app
# Prints the settings of staging app

# Specifying other apps during an operation
$ stream-cli chat get-app --app prod
# Prints the settings of prod app

```

### Options

```
  -h, --help   help for default
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli config](stream-cli_config.md)	 - Manage app configurations

