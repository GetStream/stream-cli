## stream-cli chat translate-message

Translate a message

### Synopsis

Chat messages can be translated on-demand or automatically, this
allows users speaking different languages on the same channel.

The translate endpoint returns the translated message, updates
it and sends a message.updated event to all users on the channel.


```
stream-cli chat translate-message --message-id [message-id] --language [language] --output-format [json|tree] [flags]
```

### Examples

```
# Translates a message with id 'msgid-1' to English
$ stream-cli chat translate-message --message-id msgid-1 --language en

```

### Options

```
  -h, --help                   help for translate-message
  -l, --language string        [required] Language to translate to
  -m, --message-id string      [required] Message id to translate
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

