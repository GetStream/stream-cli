## stream-cli chat list-members

List members of a channel

### Synopsis

List and paginate members of a channel using the Stream Chat API.

This command supports optional filters, offset/limit for pagination,
and sort fields (e.g. "user_id", "created_at").


```
stream-cli chat list-members --type [channel-type] --id [channel-id] [flags]
```

### Examples

```
# List first 10 members of 'red-team'
$ stream-cli chat list-members --type messaging --id red-team

# Filter members whose name includes 'tom'
$ stream-cli chat list-members --type messaging --id red-team --filter '{"name":{"$q":"tom"}}'

# Get next page of members (10 offset)
$ stream-cli chat list-members --type messaging --id red-team --offset 10 --limit 10

# Sort members by user_id ascending
$ stream-cli chat list-members --type messaging --id red-team --sort user_id:1

```

### Options

```
      --filter string          [optional] JSON string to filter members
  -h, --help                   help for list-members
  -i, --id string              [required] Channel ID
      --limit int              [optional] Pagination limit (default 10) (default 10)
      --offset int             [optional] Pagination offset (default 0)
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
      --sort string            [optional] Sorting field and direction (e.g., user_id:1 or created_at:-1)
  -t, --type string            [required] Channel type such as 'messaging' or 'livestream'
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

