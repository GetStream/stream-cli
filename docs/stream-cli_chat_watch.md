## stream-cli chat watch

Wait for an async task to complete

### Synopsis

This command waits for a specific async backend operation
to complete. Such as deleting a user or exporting a channel.


```
stream-cli chat watch [task-id] [flags]
```

### Examples

```
# Delete user and watching it complete
$ stream-cli chat delete-user --user "my-user-1"
> Successfully initiated user deletion. Task id: 7586fa0d-dc8d-4f6f-be2d-f952d0e26167

# Waiting for the task to complete
$ stream-cli chat watch 7586fa0d-dc8d-4f6f-be2d-f952d0e26167

# Providing a timeout of 80 seconds
$ stream-cli chat watch 7586fa0d-dc8d-4f6f-be2d-f952d0e26167 --timeout 80

```

### Options

```
  -h, --help          help for watch
  -t, --timeout int   [optional] Timeout in seconds. Default is 30 (default 30)
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

