# 📃 Use cases <!-- omit in toc -->

- [App configuration](#app-configuration)
- [Channel](#channel)
- [Imports](#imports)
- [GDPR](#gdpr)
- [Moderation](#moderation)

## App configuration

You can handle multiple Chat apps with the CLI.

List your configurations:
```shell
$ stream-cli config list

           Name     Access Key         Secret Key          URL
           ----     ----------         ----------          ---
(default)  test     kujk3ms96pby       **************323m  https://chat.stream-io-api.com
           prod     v5hg34n2m2nv       **************76as  https://chat.stream-io-api.com
           staging  nrnn2rmnb52u       **************242b  https://chat.stream-io-api.com
```

- [**config list** docs](./stream-cli_config_list.md)

The **default** app is used when no `--app` flag is provided to the command.

Add a new configuration:

```shell
$ stream-cli config new
? What is the name of your app? (eg. prod, staging, testing) prod
? What is your access key? v5hg34n2m2nv
? What is your access secret key? ***********************************************************
? (optional) Which base URL do you want to use for Chat? https://chat.stream-io-api.com
Application successfully added. 🚀
```
- [**config new** docs](./stream-cli_config_new.md)

From that point, you can provide `--app prod` as an argument to any command. Example:
```shell
# Create a new channel in the prod app
$ stream-cli chat create-channel -i redteam -t messaging -u joe --app prod
```

Delete a configuration:
```shell
$ stream-cli config remove prod
[prod] application successfully deleted.
```
- [**config remove** docs](./stream-cli_config_remove.md)

## Channel

All CRUD channel operations are available in the CLI.

Create a channel:
```shell
$ stream-cli chat create-channel -i redteam -t messaging -u joe
Successfully created channel [messaging:redteam2]
```
- [**create-channel** docs](./stream-cli_chat_create-channel.md)

Add members to a channel:
```shell
$ stream-cli chat add-members --type messaging --id red-team joe jill jane
Successfully added user(s) to channel
```
- [**add-members** docs](./stream-cli_chat_add-members.md)

Send a message to a channel:
```shell
$ stream-cli chat send-message --channel-type messaging --channel-id redteam --text "Hello World!" --user joe
Message successfully sent. Message id: [74c63670-f5ea-4b62-a149-98f434f321c1]
```
- [**send-message** docs](./stream-cli_chat_send-message.md)

Send a reaction:
```shell
$ stream-cli chat send-reaction --message-id 74c63670-f5ea-4b62-a149-98f434f321c1 --user user --reaction-type like
Successfully sent reaction
```

- [**send-reaction** docs](./stream-cli_chat_send-reaction.md)

List channels:
```shell
$ stream-cli chat list-channels -t messaging
< json payload >
```
- [**list-channels** docs](./stream-cli_chat_list-channels.md)


## Imports

Upload a new import:
```shell
$ stream-cli chat upload-import data.json --mode insert
```
- [**upload-import** docs](./stream-cli_chat_upload-import.md)

- [Imports how-to](./imports.md)

## GDPR

Delete users:
```shell
$ stream-cli chat delete-users joe jill
```
- [**delete-users** docs](./stream-cli_chat_delete-users.md)

Delete channel:
```shell
$ stream-cli chat delete-channel --type messaging --id redteam
Successfully initiated channel deletion. Task id: 66bbcdcd-b133-43ce-ab63-557c14d2a168

# Wait for the task to complete
$ stream-cli chat watch 66bbcdcd-b133-43ce-ab63-557c14d2a168
Waiting for async task to complete...⏳
Still loading... ⏳
Async operation completed successfully
```
- [**delete-channel** docs](./stream-cli_chat_delete-channel.md)
- [**watch** docs](./stream-cli_chat_watch.md)

## Moderation

Ban a user:
```shell
$ stream-cli chat ban-user --target-user-id mike --banned-by admin-user-2 --reason "Bad behavior"
```
- [**ban-user** docs](./stream-cli_chat_ban-user.md)

Unban a user:
```shell
$ stream-cli chat unban-user --target-user-id joe
```
- [**unban-user** docs](./stream-cli_chat_unban-user.md)

Flag a message:
```shell
$ stream-cli chat flag-message --message-id msgid-1 --user-id userid-1
Successfully flagged message.
```
- [**flag-message** docs](./stream-cli_chat_flag-message.md)

Mute a user:
```shell
$ stream-cli chat mute-user --target-user-id joe --muted-by-id admin --expiration 5
Successfully muted user.
```
- [**mute-user** docs](./stream-cli_chat_mute-user.md)

Unmute a user:
```shell
$ stream-cli chat unmute-user --target-user-id joe --unmuted-by-id admin
Successfully unmuted user.
```
- [**unmute-user** docs](./stream-cli_chat_unmute-user.md)