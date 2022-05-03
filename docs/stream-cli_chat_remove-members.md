## stream-cli chat remove-members

Remove members from a channel

```
stream-cli chat remove-members --type [channel-type] --id [channel-id] [user-id-1] [user-id-2] ... [flags]
```

### Examples

```
# Remove members joe, jill and jane from 'red-team' channel
$ stream-cli chat remove-members --type messaging --id red-team joe jill jane

```

### Options

```
  -h, --help          help for remove-members
  -i, --id string     [required] Channel id
  -t, --type string   [required] Channel type such as 'messaging' or 'livestream'
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

