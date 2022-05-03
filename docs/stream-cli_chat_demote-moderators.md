## stream-cli chat demote-moderators

Demote users from moderator role

```
stream-cli chat demote-moderators --type [channel-type] --id [channel-id] [user-id-1] [user-id-2] ... [flags]
```

### Examples

```
# Demote 4 users from moderator role
$ stream-cli chat demote-moderators --type messaging --id red-team joe mike jane jill

```

### Options

```
  -h, --help          help for demote-moderators
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

