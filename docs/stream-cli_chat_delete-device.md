## stream-cli chat delete-device

Delete a device

### Synopsis

Unregistering a device removes the device from the user
and stops further new message notifications.


```
stream-cli chat delete-device --id [device-id] --user-id [user-id] [flags]
```

### Examples

```
# Delete "my-device-id" device
$ stream-cli chat delete-device --id "my-device-id" --user-id "my-user-id"

```

### Options

```
  -h, --help             help for delete-device
  -i, --id string        [required] Device ID to delete
  -u, --user-id string   [required] ID of the user who deletes the device
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

