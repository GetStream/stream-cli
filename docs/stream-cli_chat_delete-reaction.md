## stream-cli chat delete-reaction

Delete a reaction from a message

```
stream-cli chat delete-reaction --message-id [message-id] --reaction-type [reaction-type] --user-id [user-id] [flags]
```

### Examples

```
# Delete a reaction from [08f64828-3bba-42bd-8430-c26a3634ee5c] message
$ stream-cli chat delete-reaction --message-id 08f64828-3bba-42bd-8430-c26a3634ee5c --reaction-type like --user-id 12345

```

### Options

```
  -h, --help                   help for delete-reaction
  -m, --message-id string      [required] The message id to delete the reaction from
  -r, --reaction-type string   [required] The reaction type to delete
  -u, --user-id string         [required] The user id of the user deleting the reaction
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

