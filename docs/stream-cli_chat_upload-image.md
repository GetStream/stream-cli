## stream-cli chat upload-image

Upload an image

### Synopsis

Stream supported image types are: image/bmp, image/gif, image/jpeg, image/png,
image/webp, image/heic, image/heic-sequence, image/heif, image/heif-sequence,
image/svg+xml.
You can set a more restrictive list for your application if needed.
The maximum file size is 100MB.
Stream will allow any file extension. If you want to be more restrictive
for an application, this is can be set via API or by logging into your dashboard.


```
stream-cli chat upload-image --channel-type [channel-type] --channel-id [channel-id] --user-id [user-id] --file [file] --content-type [content-type] [flags]
```

### Examples

```
# Uploads an image to 'redteam' channel of 'messaging' channel type
$ stream-cli chat upload-image --channel-type messaging --channel-id redteam --user-id "user-1" --file "./picture.png" --content-type "image/png"

```

### Options

```
  -i, --channel-id string     [required] Channel id to interact with
  -t, --channel-type string   [required] Channel type to interact with
  -f, --file string           [required] Image file path
  -h, --help                  help for upload-image
  -u, --user-id string        [required] User id
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

