## stream-cli chat create-device

Create a device

### Synopsis

Registering a device associates it with a user and tells
the push provider to send new message notifications to the device.


```
stream-cli chat create-device --id [device-id] --push-provider [firebase|apn|xiaomi|huawei] --push-provider-name [provider-name] --user-id [user-id] [flags]
```

### Examples

```
# Create a device with a firebase push provider
$ stream-cli chat create-device --id "my-device-id" --push-provider firebase --push-provider-name "my-firebase-project-id" --user-id "my-user-id"

```

### Options

```
  -h, --help                        help for create-device
  -i, --id string                   [required] Device ID
  -p, --push-provider string        [required] Push provider. Can be apn, firebase, xiaomi, huawei
  -n, --push-provider-name string   [optional] Push provider name
  -u, --user-id string              [required] User ID
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

