## stream-cli chat mute-channel

Mute a channel for a user

```
stream-cli chat mute-channel --type [channel-type] --id [channel-id] --user-id [user-id] [--expiration duration] [flags]
```

### Synopsis

Mutes a channel for a specific user. Muted channels do not trigger notifications,  
affect unread counts, or unhide themselves when new messages are added.

You can optionally set an expiration time for the mute using the `--expiration` flag,  
such as `'1h'`, `'24h'`, etc.

### Examples

```
# Mute a channel indefinitely for user 'john'
$ stream-cli chat mute-channel --type messaging --id redteam --user-id john

# Mute a channel for 6 hours
$ stream-cli chat mute-channel --type messaging --id redteam --user-id john --expiration 6h
```

### Options

```
  -t, --type string          [required] Channel type such as 'messaging'
  -i, --id string            [required] Channel ID
  -u, --user-id string       [required] User ID
      --expiration string    [optional] Expiration duration (e.g., '1h', '24h')
  -h, --help                 help for mute-channel
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications
