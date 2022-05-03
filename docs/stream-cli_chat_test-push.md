## stream-cli chat test-push

Test push notifications

```
stream-cli chat test-push --message-id [string] --apn-template [string] --firebase-template [string] --firebase-data-template [string] --skip-devices [true|false] --push-provider-name [string] --push-provider-type [string] --user-id [string] --output-format [json|tree] [flags]
```

### Examples

```
# A test push notification for a certain message id
$ stream-cli chat test-push --message-id msgid --user-id id --skip-devices true

```

### Options

```
  -h, --help                        help for test-push
      --message-id string           [optional] Message id to test
  -o, --output-format string        [optional] Output format. Can be json or tree (default "json")
      --push-provider-name string   [optional] Push provider name to use
      --push-provider-type string   [optional] Push provider type to use
      --skip-devices                [optional] Whether to notify devices
      --user-id string              [optional] User id to initiate the test
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

