## stream-cli chat list-devices

List devices

### Synopsis

Provides a list of all devices associated with a user.

```
stream-cli chat list-devices --user-id [user-id] --output-format [json|tree] [flags]
```

### Examples

```
# List devices for a user
$ stream-cli chat list-devices --user-id "my-user-id"

```

### Options

```
  -h, --help                   help for list-devices
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
  -u, --user-id string         [required] User ID
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

