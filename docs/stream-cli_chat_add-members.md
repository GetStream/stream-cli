## stream-cli chat add-members

Add members to a channel

```
stream-cli chat add-members --type [channel-type] --id [channel-id] [user-id-1] [user-id-2] ... [flags]
```

### Examples

```
# Add members joe, jill and jane to 'red-team' channel
$ stream-cli chat add-member --type messaging --id red-team joe jill jane

```

### Options

```
  -h, --help          help for add-members
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

