## stream-cli chat create-channel

Create a channel

### Synopsis

This command allows you to create a new channel. If it
exists already an error will be thrown.


```
stream-cli chat create-channel --type [channel-type] --id [channel-id] --user [user-id] [flags]
```

### Examples

```
# Create a channel with id 'redteam' of type 'messaging' by 'joe'
$ stream-cli chat create-channel --type messaging --id redteam --user joe

```

### Options

```
  -h, --help          help for create-channel
  -i, --id string     [required] Channel id
  -t, --type string   [required] Channel type such as 'messaging' or 'livestream'
  -u, --user string   [required] User id who will be considered as the creator of the channel
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

