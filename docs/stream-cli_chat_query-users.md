## stream-cli chat query-users

Query users

### Synopsis

This command allows you to search for users. The 'filter' flag is a raw JSON string,
and you can check the valid combinations in the official documentation.

https://getstream.io/chat/docs/node/query_users/?language=javascript


```
stream-cli chat query-users --filter [raw-json] --limit [int] --output-format [json|tree] [flags]
```

### Examples

```
# Query for 'user-1'. The results are shown as json.
$ stream-cli chat query-users --filter '{"id": {"$eq": "user-1"}}'

# Query for 'user-1' and 'user-2'. The results are shown as a browsable tree.
$ stream-cli chat query-users --filter '{"id": {"$in": ["user-1", "user-2"]}}' --output-format tree

```

### Options

```
  -f, --filter string          [required] Filter for users (default "{}")
  -h, --help                   help for query-users
  -l, --limit int              [optional] The number of users returned (default 10)
  -o, --output-format string   [optional] Output format. Can be json or tree (default "json")
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

