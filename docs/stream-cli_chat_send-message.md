## stream-cli chat send-message

Send a message to a channel

```
stream-cli chat send-message --channel-type [channel-type] --channel-id [channel-id] --text [text] --user [user-id] [flags]
```

### Examples

```
# Sends a message to 'redteam' channel of 'messaging' channel type
$ stream-cli chat send-message --channel-type messaging --channel-id redteam --text "Hello World!" --user "user-1"

# Sends a message to 'redteam' channel of 'livestream' channel type with an URL attachment
$ stream-cli chat send-message --channel-type livestream --channel-id redteam --attachment "https://example.com/image.png" --text "Hello World!" --user "user-1"

# You can also send a message with a local file attachment
# In this scenario, we'll upload the file first then send the message
$ stream-cli chat send-message --channel-type livestream --channel-id redteam --attachment "./image.png" --text "Hello World!" --user "user-1"

```

### Options

```
  -a, --attachment string     [optional] URL of the an attachment
  -i, --channel-id string     [required] Channel id
  -t, --channel-type string   [required] Channel type such as 'messaging' or 'livestream'
  -h, --help                  help for send-message
      --text string           [required] Text of the message
  -u, --user string           [required] User id
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

