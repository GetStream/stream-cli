## stream-cli chat create-token

Create a token

### Synopsis

Stream uses JWT (JSON Web Tokens) to authenticate chat users, enabling them to login.
Knowing whether a user is authorized to perform certain actions is
managed separately via a role based permissions system.

With this command you can generate token for a specific user that can be
used on the frontend.


```
stream-cli chat create-token --user [user-id] --expiration [epoch] --issued-at [epoch] [flags]
```

### Examples

```
# Create a JWT token for a user with id '123'. This token has no expiration.
$ stream-cli chat create-token --user 123

# Create a JWT for user 'joe' with 'exp' and 'iat' claim
$ stream-cli chat create-token --user joe --expiration 1577880000 --issued-at 1577880000

```

### Options

```
  -e, --expiration int   [optional] Expiration (exp) of the JWT in epoch timestamp
  -h, --help             help for create-token
  -i, --issued-at int    [optional] Issued at (iat) of the JWT in epoch timestamp
  -u, --user string      [required] Id of the user to create token for
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

