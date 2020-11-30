`stream config`
===============

Configure API access

* [`stream config:destroy`](#stream-configdestroy)
* [`stream config:get`](#stream-configget)
* [`stream config:set`](#stream-configset)

## `stream config:destroy`

Destroys your user configuration.

```
USAGE
  $ stream config:destroy

OPTIONS
  -f, --force  Force remove Stream configuration from cache.
```

## `stream config:get`

Outputs your user configuration.

```
USAGE
  $ stream config:get

OPTIONS
  -j, --json  Output results in JSON. When not specified, returns output in a human friendly format.
```

## `stream config:set`

Sets your user configuration.

```
USAGE
  $ stream config:set

OPTIONS
  -e, --email=email              Email for configuration.
  -j, --json                     Output results in JSON. When not specified, returns output in a human friendly format.
  -k, --key=key                  API key for configuration.
  -n, --name=name                Full name for configuration.
  -o, --timeout=timeout          Timeout for requests in ms.
  -s, --secret=secret            API secret for configuration.
  -t, --telemetry                Enable error reporting for debugging purposes.
  -u, --url=url                  API base URL for configuration.
  -v, --environment=environment  Environment to run in (production or development for token and permission checking).
```
