## stream-cli chat reactivate-user

Reactivate a user

### Synopsis

Deactivated users cannot connect to Stream Chat or send/receive messages.
This function reactivates a user.


```
stream-cli chat reactivate-user --user-id [user-id] --restore-messages [true|false] [flags]
```

### Examples

```
# Reactivate the user 'joe'
$ stream-cli chat reactivate-user --user-id joe --restore-messages true

```

### Options

```
  -h, --help               help for reactivate-user
      --restore-messages   [optional] Restore messages for the user
  -u, --user-id string     [required] ID of the user to reactivate
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

