## stream-cli chat deactivate-user



### Synopsis

Deactivated users cannot connect to Stream Chat or send/receive messages.
Deactivated users can be re-activated with the 'reactivate-user' command.


```
stream-cli chat deactivate-user --user-id [user-id] --mark-messages-deleted [true|false] [flags]
```

### Examples

```
# Deactivate the user 'joe'
$ stream-cli chat deactivate-user --user-id joe --mark-messages-deleted true

```

### Options

```
  -h, --help                    help for deactivate-user
      --mark-messages-deleted   [optional] Mark all messages from the user as deleted
  -u, --user-id string          [required] ID of the user to deactivate
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

