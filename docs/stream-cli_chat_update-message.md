## stream-cli chat update-message

Update an existing message

### Synopsis

Update a message by providing the message ID, user ID, and new message text.
This fully overwrites the message content while preserving metadata.


```
stream-cli chat update-message --message-id [id] --user [user-id] --text [text] [flags]
```

### Examples

```
# Update a message by ID
$ stream-cli chat update-message --message-id msgid-123 --user user123 --text "Updated message text"

```

### Options

```
  -h, --help                help for update-message
      --message-id string   [required] ID of the message to update
      --text string         [required] New message text
      --user string         [required] User ID performing the update
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

