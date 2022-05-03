## stream-cli chat ban-user

Ban a user

### Synopsis

Users can be banned from an app entirely.
When a user is banned, they will not be allowed to post messages until
the ban is removed or expired but will be able to connect to Chat
and to channels as before.

Channel watchers cannot be banned.


```
stream-cli chat ban-user --target-user-id [user-id] --banned-by-id [user-id] --reason [reason] --expiration [expiration-in-minutes] [flags]
```

### Examples

```
# 'admin-user-1' bans user 'joe'
$ stream-cli chat ban-user --target-user-id joe --banned-by admin-user-1

# 'admin-user-2' bans user 'mike' with a reason
$ stream-cli chat ban-user --target-user-id mike --banned-by admin-user-2 --reason "Bad behavior"

# 'admin-user-3' bans user 'jill' with a reason for 1 hour
$ stream-cli chat ban-user --target-user-id jill --banned-by admin-user-3 --expiration 60

```

### Options

```
  -b, --banned-by-id string     [required] ID of the user who is performing the ban
  -e, --expiration int          [optional] Number of minutes until the ban expires. Defaults to forever.
  -h, --help                    help for ban-user
  -r, --reason string           [optional] Reason for the ban
  -t, --target-user-id string   [required] ID of the user to ban
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

