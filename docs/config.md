`stream config`
===============

manage configuration variables

* [`stream config:destroy`](#stream-configdestroy)
* [`stream config:get`](#stream-configget)
* [`stream config:set`](#stream-configset)

## `stream config:destroy`

```
USAGE
  $ stream config:destroy

OPTIONS
  -f, --force  Force remove Stream configuration from cache.
```

_See code: [src/commands/config/destroy.js](https://github.com/getstream/stream-cli/blob/v0.0.2/src/commands/config/destroy.js)_

## `stream config:get`

```
USAGE
  $ stream config:get

OPTIONS
  -j, --json  Output results in JSON. When not specified, returns output in a human friendly format.
```

_See code: [src/commands/config/get.js](https://github.com/getstream/stream-cli/blob/v0.0.2/src/commands/config/get.js)_

## `stream config:set`

```
USAGE
  $ stream config:set

OPTIONS
  -e, --email=email    Email for configuration.
  -j, --json           Output results in JSON. When not specified, returns output in a human friendly format.
  -k, --key=key        API key for configuration.
  -n, --name=name      Full name for configuration.
  -s, --secret=secret  API secret for configuration.
```

_See code: [src/commands/config/set.js](https://github.com/getstream/stream-cli/blob/v0.0.2/src/commands/config/set.js)_
