## stream-cli chat upload-file

Upload a file

### Synopsis

Stream will not block any file types from uploading, however, different
clients may handle different types differently or not at all.
You can set a more restrictive list for your application if needed.
The maximum file size is 100MB.
Stream will allow any file extension. If you want to be more restrictive
for an application, this is can be set via API or by logging into your dashboard.


```
stream-cli chat upload-file --channel-type [channel-type] --channel-id [channel-id] --user-id [user-id] --file [file] [flags]
```

### Examples

```
# Uploads a file to 'redteam' channel of 'messaging' channel type
$ stream-cli chat upload-file --channel-type messaging --channel-id redteam --user-id "user-1" --file "./snippet.txt"

```

### Options

```
  -i, --channel-id string     [required] Channel id to interact with
  -t, --channel-type string   [required] Channel type to interact with
  -f, --file string           [required] File path
  -h, --help                  help for upload-file
  -u, --user-id string        [required] User id
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

