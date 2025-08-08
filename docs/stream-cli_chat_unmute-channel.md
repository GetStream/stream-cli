## stream-cli chat unmute-channel

Unmute a channel for a user

```
stream-cli chat unmute-channel --type [channel-type] --id [channel-id] --user-id [user-id] [flags]
```

### Synopsis

Unmutes a previously muted channel for a specific user.

### Examples

```
# Unmute the 'redteam' channel for user 'john'
$ stream-cli chat unmute-channel --type messaging --id redteam --user-id john
```

### Options

```
  -t, --type string        [required] Channel type such as 'messaging'
  -i, --id string          [required] Channel ID
  -u, --user-id string     [required] User ID
  -h, --help               help for unmute-channel
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications
