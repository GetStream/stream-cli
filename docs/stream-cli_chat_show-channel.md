## stream-cli chat show-channel

Show a channel

### Synopsis

Hiding a channel will remove it from query channel requests for that
user until a new message is added.
As opposed to this, showing a channel will add it to query channel requests for that user.


```
stream-cli chat show-channel --type [channel-type] --id [channel-id] --user-id [user-id] [flags]
```

### Examples

```
# Show a 'red-team' channel for user 'joe'
$ stream-cli chat show-channel --type messaging --id red-team --user-id joe

```

### Options

```
  -h, --help             help for show-channel
  -i, --id string        [required] Channel id
  -t, --type string      [required] Channel type such as 'messaging' or 'livestream'
  -u, --user-id string   [required] User id to show the channel to
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

