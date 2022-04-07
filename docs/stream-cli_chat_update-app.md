## stream-cli chat update-app

Update application settings

### Synopsis

Update the application settings.

Application level settings allow you to configure settings that
impact all the channel types in your app.

See https://getstream.io/chat/docs/rest/#settings-updateapp for
the available JSON options.


```
stream-cli chat update-app --properties [raw-json-update-properties] [flags]
```

### Examples

```
# Enable multi-tenant and update permission version to v2
$ stream-cli chat update-app --properties '{"multi_tenant_enabled": true, "permission_version": "v2"}'

```

### Options

```
  -h, --help                help for update-app
  -p, --properties string   [required] Raw json properties to update
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

