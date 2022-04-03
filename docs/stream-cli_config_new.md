## stream-cli config new

Add a new application

### Synopsis

Add a new application which can be used for further operations

```
stream-cli config new [flags]
```

### Examples

```
# Add a new application to the CLI
$ stream-cli config new
? What is the name of your app? (eg. prod, staging, testing) testing
? What is your access key? abcd1234efgh456
? What is your access secret key? ***********************************
? (optional) Which base URL do you want to use for Chat? https://chat.stream-io-api.com

Application successfully added. 🚀

```

### Options

```
  -h, --help   help for new
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli config](stream-cli_config.md)	 - Manage app configurations

