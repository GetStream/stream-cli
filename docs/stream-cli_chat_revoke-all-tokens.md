## stream-cli chat revoke-all-tokens

Revoke all tokens

### Synopsis

This command revokes ALL tokens for all users of an application.
This should be used with caution as it will expire every user’s token,
regardless of whether the token has an iat claim.


```
stream-cli chat revoke-all-tokens --before [epoch] [flags]
```

### Examples

```
# Revoke all tokens for the default app, from now
$ stream-cli chat revoke-all-tokens

# Revoke all tokens for the test app, before 2019-01-01
$ stream-cli chat revoke-all-tokens --before 1546300800 --app test

```

### Options

```
  -b, --before int   [optional] The epoch timestamp before which tokens should be revoked. Defaults to now.
  -h, --help         help for revoke-all-tokens
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

