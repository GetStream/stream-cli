## stream-cli chat send-reaction

Send a reaction to a message

### Synopsis

Stream Chat has built-in support for user Reactions. Common examples are
likes, comments, loves, etc. Reactions can be customized so that you
are able to use any type of reaction your application requires.


```
stream-cli chat send-reaction --message-id [message-id] --user-id [user-id] --reaction-type [reaction-type] [flags]
```

### Examples

```
# Send a reaction to a [08f64828-3bba-42bd-8430-c26a3634ee5c] message
$ stream-cli chat send-reaction --message-id 08f64828-3bba-42bd-8430-c26a3634ee5c --user-id 12345 --reaction-type like

```

### Options

```
  -h, --help                   help for send-reaction
  -m, --message-id string      [required] The message id to send the reaction to
  -r, --reaction-type string   [required] The reaction type to send
  -u, --user-id string         [required] The user id of the user sending the reaction
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

