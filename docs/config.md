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

## `stream config:get`

```
USAGE
  $ stream config:get

OPTIONS
  -j, --json  Output results in JSON. When not specified, returns output in a human friendly format.
```

## `stream config:set`

```
USAGE
  $ stream config:set

OPTIONS
  -e, --email=email    Email for configuration.
  -j, --json           Output results in JSON. When not specified, returns output in a human friendly format.
  -k, --key=key        API key for configuration.
  -m, --mode=mode      Environment to run in (production or development for token and permission checking).
  -n, --name=name      Full name for configuration.
  -s, --secret=secret  API secret for configuration.
  -u, --url=url        API base URL for configuration.
```
