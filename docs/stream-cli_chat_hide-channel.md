## stream-cli chat hide-channel

Hide a channel

### Synopsis

Hiding a channel will remove it from query channel requests for that
user until a new message is added. Please keep in mind that hiding a channel
is only available to members of that channel.
You can still retrieve the list of hidden channels using the { "hidden" : true } query parameter.


```
stream-cli chat hide-channel --type [channel-type] --id [channel-id] --user-id [user-id] [flags]
```

### Examples

```
# Hide a 'red-team' channel for user 'joe'
$ stream-cli chat hide-channel --type messaging --id red-team --user-id joe

```

### Options

```
  -h, --help             help for hide-channel
  -i, --id string        [required] Channel id
  -t, --type string      [required] Channel type such as 'messaging' or 'livestream'
  -u, --user-id string   [required] User id to hide the channel to
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

