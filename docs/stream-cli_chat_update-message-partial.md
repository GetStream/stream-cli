## stream-cli chat update-message-partial

Partially update a message

### Synopsis

A partial update can be used to set and unset specific fields when it
is necessary to retain additional data fields on the object. AKA a patch style update.


```
stream-cli chat update-message-partial --message-id [message-id] --user [user-id] --set [key-value-pairs] --unset [property-names] [flags]
```

### Examples

```
# Partially updates a message with id 'msgid-1'. Updates a custom field and removes the silent flag.
$ stream-cli chat update-message-partial -message-id msgid-1 --set importance=low --unset silent

```

### Options

```
  -h, --help                 help for update-message-partial
  -m, --message-id string    [required] Message id
  -s, --set stringToString   [optional] Comma-separated key-value pairs to set (default [])
      --unset string         [optional] Comma separated list of properties to unset
  -u, --user string          [required] User id
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

