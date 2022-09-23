## stream-cli chat update-user-partial

Partially update a user

### Synopsis

Updates an existing user. The 'set' property is a comma separated list of key value pairs.
The 'unset' property is a comma separated list of property names.


```
stream-cli chat update-user-partial --user [user-id] --set [raw-json] --unset [property-names] [flags]
```

### Examples

```
# Set a user's role to 'admin' and set 'age' to 21. At the same time, remove 'haircolor' and 'height'.
$ stream-cli chat update-user-partial --user-id my-user-1 --set '{"role":"admin","age":21}' --unset haircolor,height

```

### Options

```
  -h, --help             help for update-user-partial
  -s, --set string       [optional] Raw JSON of key-value pairs to set
  -u, --unset string     [optional] Comma separated list of properties to unset
  -i, --user-id string   [required] Channel ID
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

