## stream-cli chat flag-message

Flag a message

### Synopsis

Any user is allowed to flag a message. This triggers the message.flagged webhook event
and adds the message to the inbox of your Stream Dashboard Chat Moderation view.


```
stream-cli chat flag-message --message-id [message-id] --user-id [user-id] [flags]
```

### Examples

```
# Flags a message with id 'msgid-1' by 'userid-1'
$ stream-cli chat flag-message --message-id msgid-1 --user-id userid-1

```

### Options

```
  -h, --help                help for flag-message
  -m, --message-id string   [required] Message id to flag
  -u, --user-id string      [required] ID of the user who flagged the message
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

