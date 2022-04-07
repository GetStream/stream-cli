## stream-cli chat delete-channel

Delete a channel

### Synopsis

This command allows you to delete a channel. This operation is asynchronous
in the backend so a task id is returned. You need to use the watch
commnand to poll the results.


```
stream-cli chat delete-channel --type [channel-type] --id [channel-id] [flags]
```

### Examples

```
# Delete a channel with id 'redteam' of type 'messaging'
$ stream-cli chat delete-channel --type messaging --id redteam
> Successfully initiated channel deletion. Task id: 66bbcdcd-b133-43ce-ab63-557c14d2a168

# Wait for the task to complete
$ stream-cli chat watch 66bbcdcd-b133-43ce-ab63-557c14d2a168

```

### Options

```
      --hard          [optional] Channel will be hard deleted. This action is irrevocable.
  -h, --help          help for delete-channel
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

