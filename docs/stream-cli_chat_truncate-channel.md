## stream-cli chat truncate-channel

Truncate a channel by removing all messages but keeping the channel metadata and members.

```
stream-cli chat truncate-channel --type [channel-type] --id [channel-id] [flags]
```

### Examples

```
# Truncate messages in 'general' channel of type messaging
$ stream-cli chat truncate-channel --type messaging --id general

# Truncate with hard delete and system message
$ stream-cli chat truncate-channel --type messaging --id general --hard --message "Channel reset" --message-user-id system-user
```

### Options

```
  -h, --help                   help for truncate-channel
  -i, --id string              [required] Channel ID
      --user-id string         [optional] User ID who performs the truncation
      --message string         [optional] System message to include in truncation (requires --message-user-id)
      --message-user-id string [optional] User id for the message to include in truncation (required if --message is set)
      --hard                   [optional] Permanently delete messages instead of hiding them
      --skip-push               [optional] Skip push notifications
  -t, --type string            [required] Channel type such as 'messaging'
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

