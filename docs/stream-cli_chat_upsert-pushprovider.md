## stream-cli chat upsert-pushprovider

Create or updates a push provider

### Synopsis


			The "--properties" parameter expects a raw json string that can be
			unmarshalled into a stream_chat.PushProvider object on the Go SDK side.
			See the example section.
			Available properties:
			type
			name
			description
			disabled_at
			disabled_reason
					
			apn_auth_key
			apn_key_id
			apn_team_id
			apn_topic

			firebase_notification_template
			firebase_apn_template
			firebase_credentials
					
			huawei_app_id
			huawei_app_secret
					
			xiaomi_package_name
			xiaomi_app_secret
		

```
stream-cli chat upsert-pushprovider --properties [raw-json] [flags]
```

### Examples

```
# Setting up an APN push provider
$ stream-cli chat upsert-pushprovider --properties "{'type': 'apn', 'name': 'staging', 'apn_auth_key': 'key', 'apn_key_id': 'id', 'apn_topic': 'topic', 'apn_team_id': 'id'}"

# Setting up a Firebase push provider
$ stream-cli chat upsert-pushprovider --properties "{'type': 'firebase', 'name': 'staging', 'firebase_credentials': 'credentials'}"

# Setting up a Huawei push provider
$ stream-cli chat upsert-pushprovider --properties "{'type': 'huawei', 'name': 'staging', 'huawei_app_id': 'id', 'huawei_app_secret': 'secret'}"

# Setting up a Xiaomi push provider
$ stream-cli chat upsert-pushprovider --properties "{'type': 'xiaomi', 'name': 'staging', 'xiaomi_package_name': 'name', 'xiaomi_app_secret': 'secret'}"

```

### Options

```
  -h, --help                help for upsert-pushprovider
  -p, --properties string   [required] Raw json properties to send to the backend
```

### Options inherited from parent commands

```
      --app string      [optional] Application name to use as it's defined in the configuration file
      --config string   [optional] Explicit config file path
```

### SEE ALSO

* [stream-cli chat](stream-cli_chat.md)	 - Allows you to interact with your Chat applications

