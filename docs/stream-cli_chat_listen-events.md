## stream-cli chat listen-events

Listen to events

### Synopsis

The command opens a WebSocket connection to the backend in the name of the user
and prints the received events to the standard output.
Press Ctrl+C to exit.


```
stream-cli chat listen-events --user-id [user-id] --timeout [number] [flags]
```

### Examples

```
# Listen to events for user with id 'my-user-1'
$ stream-cli chat listen-events --user-id my-user-1

# Listen to events for user with id 'my-user-2' and keeping the connection open for 120 seconds
$ stream-cli chat listen-events --user-id my-user-1 --timeout 120

```

### Options

```
  -h, --help                   help for listen-events
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
  -t, --timeout int32          [optional] For how many seconds do we keep the connection alive. Default is 60 seconds, max is 300. (default 60)
  -u, --user-id string         [required] User ID
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

