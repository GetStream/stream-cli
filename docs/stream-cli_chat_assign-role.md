## stream-cli chat assign-role

Assign a role to a user

```
stream-cli chat assign-role --type [channel-type] --id [channel-id] --user-id [user-id] --role [channel-role-name] [flags]
```

### Examples

```
# Assign 'channel_moderator' role to user 'joe'
$ stream-cli chat assign-role --type messaging --id red-team --user-id joe --role channel_moderator

```

### Options

```
  -h, --help             help for assign-role
  -i, --id string        [required] Channel id
  -r, --role string      [required] Channel role name to assign
  -t, --type string      [required] Channel type such as 'messaging' or 'livestream'
  -u, --user-id string   [required] User id to assign a role to
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

