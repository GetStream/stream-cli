## stream-cli chat upsert-user

Upsert a user

### Synopsis

This command inserts a new or updates an existing user.
Stream Users require only an id to be created.
Any user present in the payload will have its data replaced with the new version.


```
stream-cli chat upsert-user --properties [raw-json] [flags]
```

### Examples

```
# Create a new user with id 'my-user-1'
$ stream-cli chat upsert-user --properties "{\"id\":\"my-user-1\"}"

Check the Go SDK's 'User' struct for the properties that you can use here.

```

### Options

```
  -h, --help                help for upsert-user
  -p, --properties string   [required] Raw JSON properties of the user
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

