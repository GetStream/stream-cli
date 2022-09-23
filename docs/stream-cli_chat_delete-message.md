## stream-cli chat delete-message

Delete a message

### Synopsis

You can delete a message by calling DeleteMessage and including a message
with an ID. Messages can be soft deleted or hard deleted. Unless specified
via the hard parameter, messages are soft deleted. Be aware that deleting
a message doesn't delete its attachments.


```
stream-cli chat delete-message [message-id] [flags]
```

### Examples

```
# Soft deletes a message with id 'msgid-1'
$ stream-cli chat delete-message msgid-1

# Hard deletes a message with id 'msgid-2'
$ stream-cli chat delete-message msgid-2 --hard

```

### Options

```
  -H, --hard   [optional] Hard delete message. Default is false
  -h, --help   help for delete-message
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

