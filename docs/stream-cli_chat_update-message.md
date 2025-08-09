## stream-cli chat update-message

Update an existing message

```
stream-cli chat update-message --message-id [id] --user [user-id] --text [text]
```

### Examples

```
# Update a message by ID
$ stream-cli chat update-message --message-id msgid-123 --user user123 --text "Updated message text"
```

### Options

```
  -h, --help             help for update-message
      --message-id id    [required] Message ID to update
      --user string      [required] User ID performing the update
      --text string      [required] New message text
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications
