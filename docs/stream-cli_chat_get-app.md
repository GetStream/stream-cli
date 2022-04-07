## stream-cli chat get-app

Get application settings

### Synopsis

Get the application settings.

Application level settings allow you to configure settings that
impact all the channel types in your app.


```
stream-cli chat get-app --output-format [json|tree] [flags]
```

### Examples

```
# Print the application settings in json format (default format)
$ stream-cli chat get-app

# Print the application settings in a browsable tree
$ stream-cli chat get-app --output-format tree

# Print the application settings for another application
$ stream-cli chat get-app --app testenvironment

# Note:
# Use this command to list all the available Stream applications
$ stream-cli config list

```

### Options

```
  -h, --help                   help for get-app
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

