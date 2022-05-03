## stream-cli chat flag-user

Flag a user

```
stream-cli chat flag-user --user-id [user-id] --flagged-by-id [user-id] [flags]
```

### Examples

```
# Flag the user 'joe'
$ stream-cli chat flag-user --user-id joe --flagged-by-id admin

```

### Options

```
  -b, --flagged-by-id string   [required] ID of the user who flagged the user
  -h, --help                   help for flag-user
  -u, --user-id string         [required] ID of the user to flag
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

