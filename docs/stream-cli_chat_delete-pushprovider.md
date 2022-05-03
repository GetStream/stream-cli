## stream-cli chat delete-pushprovider

Delete a push provider

```
stream-cli chat delete-pushprovider --push-provider-type [type] --push-provider-name [name] [flags]
```

### Examples

```
# Delete an APN push provider
$ stream-cli chat delete-pushprovider --push-provider-type apn --push-provider-name staging

```

### Options

```
  -h, --help                        help for delete-pushprovider
  -n, --push-provider-name string   [required] Push provider name
  -t, --push-provider-type string   [required] Push provider type
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

