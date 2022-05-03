## stream-cli chat mute-user

Mute a user

### Synopsis

Any user is allowed to mute another user. Mutes are stored at the user
level and returned with the rest of the user information when connectUser is called.
A user will be muted until the user is unmuted or the mute is expired.


```
stream-cli chat mute-user --target-user-id [user-id] --muted-by-id [user-id] --expiration [minutes] [flags]
```

### Examples

```
# Mute the user 'joe' for 5 minutes
$ stream-cli chat mute-user --target-user-id joe --muted-by-id admin --expiration 5

```

### Options

```
  -e, --expiration int          [optional] Number of minutes until the mute expires
  -h, --help                    help for mute-user
  -b, --muted-by-id string      [required] ID of the user who muted the user
  -t, --target-user-id string   [required] ID of the user to mute
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

