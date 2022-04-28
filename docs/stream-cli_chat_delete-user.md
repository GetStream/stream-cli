## stream-cli chat delete-user

Delete a user

### Synopsis

This command deletes a user. If not flags are provided, user and messages will be soft deleted.

There are 3 additional options that you can provide:

--hard-delete: If set to true, hard deletes everything related to this user, channels, messages and everything related to it.
--mark-messages-deleted: If set to true, hard deletes all messages related to this user.
--delete-conversations: If set to true, hard deletes all conversations related to this user.

User deletion is an async operation in the backend.
Once it succeeded, you'll need to use the 'watch' command to see the async task's result.


```
stream-cli chat delete-user --user [user-id] --hard-delete [true|false] --mark-messages-deleted [true|false] --delete-conversations [true|false] [flags]
```

### Examples

```
# Soft delete a user with id 'my-user-1'
$ stream-cli chat delete-user --user my-user-1

# Hard delete a user with id 'my-user-2'
$ stream-cli chat delete-user --user my-user-2 --hard-delete
> Successfully initiated user deletion. Task id: 8d011daa-cbcd-4cba-ad16-701de599873a

# Watch the async task's result
$ stream-cli chat watch 8d011daa-cbcd-4cba-ad16-701de599873a
> Async operation completed successfully.

```

### Options

```
      --delete-conversations    [optional] Hard delete all conversations related to the user
      --hard-delete             [optional] Hard delete everything related to this user
  -h, --help                    help for delete-user
      --mark-messages-deleted   [optional] Hard delete all messages related to the user
  -u, --user string             [required] Id of the user to delete
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

