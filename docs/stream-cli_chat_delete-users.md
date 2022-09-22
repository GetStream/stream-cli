## stream-cli chat delete-users

Delete multiple users

### Synopsis

You can delete up to 100 users and optionally all of their channels and
messages using this method.  First the users are marked deleted synchronously
so the user will not be directly visible in the API.  Then the process deletes
the user and related objects asynchronously by scheduling a task to be handle
by the task worker.

The delete users response contain a task ID which can be polled using the
get task endpoint to check the status of the deletions.

Note: when deleting a user with hard delete, it also required hard deletion
for messages and conversations as well!


```
stream-cli chat delete-users --new-channel-user-id [user-id] --hard-delete-users [true|false] --hard-delete-messages [true|false] --hard-delete-conversations [user-id1] [user-id2] ... [flags]
```

### Examples

```
# Soft delete users with ids 'my-user-1' and 'my-user-2'
$ stream-cli chat delete-users my-user-1 my-user-2
> Successfully initiated user deletion. Task id: bf1c2d1b-04d6-4e67-873c-5b3ade478b0a
# Waiting for it to succeed
$ stream-cli chat watch bf1c2d1b-04d6-4e67-873c-5b3ade478b0a
> Async operation completed successfully.

# Hard delete users with ids 'my-user-3' and 'my-user-4'
$ stream-cli chat delete-users --hard-delete-users --hard-delete-messages --hard-delete-conversations my-user-3 my-user-4
> Successfully initiated user deletion. Task id: 71516d9a-0764-4aa8-b017-8d2a99748e16

# Waiting for it to succeed
$ stream-cli chat watch 71516d9a-0764-4aa8-b017-8d2a99748e16
> Async operation completed successfully.

```

### Options

```
      --hard-delete-conversations     [optional] Hard delete all conversations related to the users
      --hard-delete-messages          [optional] Hard delete all messages related to the users
      --hard-delete-users             [optional] Hard delete everything related to the users
  -h, --help                          help for delete-users
      --new-channel-owner-id string   [optional] Channels owned by hard-deleted users will be transferred to this userID
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

