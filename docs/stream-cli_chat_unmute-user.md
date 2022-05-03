## stream-cli chat unmute-user

Unmute a user

```
stream-cli chat unmute-user --target-user-id [user-id] --unmuted-by-id [user-id] [flags]
```

### Examples

```
# Unmute the user 'joe'
$ stream-cli chat unmute-user --target-user-id joe --unmuted-by-id admin

```

### Options

```
  -h, --help                    help for unmute-user
  -t, --target-user-id string   [required] ID of the user to unmute
  -b, --unmuted-by-id string    [required] ID of the user who unmuted the user
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

