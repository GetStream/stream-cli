## stream-cli chat revoke-token

Revoke a token

### Synopsis

Revokes a token for a single user. All requests will be rejected that
were issued before the given epoch timestamp.


```
stream-cli chat revoke-token --user [user-id] --before [epoch] [flags]
```

### Examples

```
# Revoke token for user 'joe' before today's date (default date)
$ stream-cli revoke-token --user joe

# Revoke token for user 'mike' before 2019-01-01
$ stream-cli revoke-token --user mike --before 1546300800

```

### Options

```
  -b, --before int    [optional] The epoch timestamp before which tokens should be revoked. Defaults to now.
  -h, --help          help for revoke-token
  -u, --user string   [required] Id of the user to revoke token for
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

