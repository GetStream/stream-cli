`stream chat`
=============

configure and manage all things related to chat

* [`stream chat:channel:create`](#stream-chatchannelcreate)
* [`stream chat:channel:get`](#stream-chatchannelget)
* [`stream chat:channel:list`](#stream-chatchannellist)
* [`stream chat:channel:query`](#stream-chatchannelquery)
* [`stream chat:channel:update`](#stream-chatchannelupdate)
* [`stream chat:log`](#stream-chatlog)
* [`stream chat:message:create`](#stream-chatmessagecreate)
* [`stream chat:message:flag`](#stream-chatmessageflag)
* [`stream chat:message:list`](#stream-chatmessagelist)
* [`stream chat:message:remove`](#stream-chatmessageremove)
* [`stream chat:message:update`](#stream-chatmessageupdate)
* [`stream chat:push:apn`](#stream-chatpushapn)
* [`stream chat:push:firebase`](#stream-chatpushfirebase)
* [`stream chat:push:get`](#stream-chatpushget)
* [`stream chat:push:webhook`](#stream-chatpushwebhook)
* [`stream chat:reaction:create`](#stream-chatreactioncreate)
* [`stream chat:reaction:remove`](#stream-chatreactionremove)
* [`stream chat:user:ban`](#stream-chatuserban)
* [`stream chat:user:create`](#stream-chatusercreate)
* [`stream chat:user:flag`](#stream-chatuserflag)
* [`stream chat:user:mute`](#stream-chatusermute)
* [`stream chat:user:remove`](#stream-chatuserremove)

## `stream chat:channel:create`

```
USAGE
  $ stream chat:channel:create

OPTIONS
  -c, --channel=channel                                 [default: f30926c7-1788-43c0-a6af-564c16b39ef6] A unique ID for
                                                        the channel you wish to create.

  -d, --data=data                                       Additional data as JSON.

  -i, --image=image                                     URL to channel image.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -n, --name=name                                       Name of the channel room.

  -t, --type=livestream|messaging|gaming|commerce|team  Type of channel.
```

_See code: [src/commands/chat/channel/create.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/channel/create.js)_

## `stream chat:channel:get`

```
USAGE
  $ stream chat:channel:get

OPTIONS
  -c, --channel=channel                                 The channel ID you wish to retrieve.
  -t, --type=livestream|messaging|gaming|commerce|team  Type of channel.
```

_See code: [src/commands/chat/channel/get.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/channel/get.js)_

## `stream chat:channel:list`

```
USAGE
  $ stream chat:channel:list
```

_See code: [src/commands/chat/channel/list.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/channel/list.js)_

## `stream chat:channel:query`

```
USAGE
  $ stream chat:channel:query

OPTIONS
  -c, --channel=channel                                 The unique indentifier for the channel you want to query.
  -f, --filter=filter                                   Filters to apply to the query.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -s, --sort=sort                                       Sort to apply to the query.

  -t, --type=livestream|messaging|gaming|commerce|team  Type of channel.
```

_See code: [src/commands/chat/channel/query.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/channel/query.js)_

## `stream chat:channel:update`

```
USAGE
  $ stream chat:channel:update

OPTIONS
  -c, --channel=channel                                 The ID of the channel you wish to update.
  -d, --description=description                         Description for the channel.
  -i, --image=image                                     URL to the channel image.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -n, --name=name                                       Name of the channel room.

  -r, --reason=reason                                   Reason for changing channel.

  -t, --type=livestream|messaging|gaming|commerce|team  Type of channel.
```

_See code: [src/commands/chat/channel/update.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/channel/update.js)_

## `stream chat:log`

```
USAGE
  $ stream chat:log

OPTIONS
  -c, --channel=channel
      The channel ID you wish to log.

  -e, 
  --event=all|user.status.changed|user.watching.start|user.watching.stop|user.updated|typing.start|typing.stop|message.n
  ew|message.updated|message.deleted|message.seen|message.reaction|member.added|member.removed|channel.updated|health.ch
  eck|connection.changed|connection.recovered
      The type of event you want to listen on.

  -j, --json
      Output results in JSON. When not specified, returns output in a human friendly format.

  -t, --type=livestream|messaging|gaming|commerce|team
      The type of channel.
```

_See code: [src/commands/chat/log/index.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/log/index.js)_

## `stream chat:message:create`

```
USAGE
  $ stream chat:message:create

OPTIONS
  -c, --channel=channel                                 The ID of the channel that you would like to send a message to.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -m, --message=message                                 The message you would like to send as plaintext.

  -n, --image=image                                     Absolute URL for an avatar of the user sending the message.

  -n, --name=name                                       The name of the user sending the message.

  -t, --type=livestream|messaging|gaming|commerce|team  The type of channel.

  -u, --user=user                                       The ID of the user sending the message.
```

_See code: [src/commands/chat/message/create.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/message/create.js)_

## `stream chat:message:flag`

```
USAGE
  $ stream chat:message:flag

OPTIONS
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message  The unique identifier of the message you want to flag.
```

_See code: [src/commands/chat/message/flag.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/message/flag.js)_

## `stream chat:message:list`

```
USAGE
  $ stream chat:message:list

OPTIONS
  -c, --channel=channel                                 The ID of the channel that you would like to send a message to.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -t, --type=livestream|messaging|gaming|commerce|team  The type of channel.
```

_See code: [src/commands/chat/message/list.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/message/list.js)_

## `stream chat:message:remove`

```
USAGE
  $ stream chat:message:remove

OPTIONS
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message  The unique identifier of the message you would like to remove.
```

_See code: [src/commands/chat/message/remove.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/message/remove.js)_

## `stream chat:message:update`

```
USAGE
  $ stream chat:message:update

OPTIONS
  -a, --attachments=attachments  A JSON payload of attachments to send along with a message.
  -j, --json                     Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message          The unique identifier for the message.
  -t, --text=text                The message you would like to send as text.
```

_See code: [src/commands/chat/message/update.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/message/update.js)_

## `stream chat:push:apn`

```
USAGE
  $ stream chat:push:apn

OPTIONS
  -a, --auth_key=auth_key                            Absolute path to .p8 auth key.
  -b, --bundle_id=bundle_id                          Bundle identifier (e.g. com.apple.test).
  -c, --p12_cert=p12_cert                            Absolute path to .p12 file.
  -d, --development                                  Use development endpoint

  -j, --json                                         Output results in JSON. When not specified, returns output in a
                                                     human friendly format.

  -k, --key_id=key_id                                Key ID.

  -n, --notification_template=notification_template  JSON template for notifications.

  -t, --team_id=team_id                              Team ID.

  --disable                                          Disable APN push notifications and clear config.
```

_See code: [src/commands/chat/push/apn.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/push/apn.js)_

## `stream chat:push:firebase`

```
USAGE
  $ stream chat:push:firebase

OPTIONS
  -j, --json                                         Output results in JSON. When not specified, returns output in a
                                                     human friendly format.

  -k, --key=key                                      Server key for Firebase.

  -n, --notification_template=notification_template  JSON notification template.

  --disable                                          Disable Firebase push notifications and clear config.
```

_See code: [src/commands/chat/push/firebase.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/push/firebase.js)_

## `stream chat:push:get`

```
USAGE
  $ stream chat:push:get

OPTIONS
  -j, --json  Output results in JSON. When not specified, returns output in a human friendly format.
```

_See code: [src/commands/chat/push/get.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/push/get.js)_

## `stream chat:push:webhook`

```
USAGE
  $ stream chat:push:webhook

OPTIONS
  -j, --json     Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --url=url  Fully qualified URL for webhook support.
```

_See code: [src/commands/chat/push/webhook.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/push/webhook.js)_

## `stream chat:reaction:create`

```
USAGE
  $ stream chat:reaction:create

OPTIONS
  -c, --channel=channel                                 The unique identifier for the channel.
  -c, --message=message                                 The unique identifier for the message.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -r, --reaction=reaction                               A reaction for the message (e.g. love).

  -t, --type=livestream|messaging|gaming|commerce|team  The type of channel.
```

_See code: [src/commands/chat/reaction/create.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/reaction/create.js)_

## `stream chat:reaction:remove`

```
USAGE
  $ stream chat:reaction:remove

OPTIONS
  -c, --channel=channel                                 The unique identifier for the channel.
  -c, --message=message                                 The unique identifier for the message.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -r, --reaction=reaction                               The unique identifier for the reaction.

  -t, --type=livestream|messaging|gaming|commerce|team  The type of channel.
```

_See code: [src/commands/chat/reaction/remove.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/reaction/remove.js)_

## `stream chat:user:ban`

```
USAGE
  $ stream chat:user:ban

OPTIONS
  -d, --duration=duration  [default: 60] Duration of timeout in minutes.
  -j, --json               Output results in JSON. When not specified, returns output in a human friendly format.
  -r, --reason=reason      A reason for adding a timeout.
  -u, --type=type          Type of ban to perform (e.g. global or channel).
  -u, --user=user          The unique identifier of the user to ban.
```

_See code: [src/commands/chat/user/ban.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/user/ban.js)_

## `stream chat:user:create`

```
USAGE
  $ stream chat:user:create

OPTIONS
  -c, --channel=channel                                                                Channel identifier.

  -j, --json                                                                           Output results in JSON. When not
                                                                                       specified, returns output in a
                                                                                       human friendly format.

  -r, --role=admin|guest|channel_moderator|channel_member|channel_owner|message_owner  The role to assign to the user.

  -t, --type=livestream|messaging|gaming|commerce|team                                 The type of channel.

  -u, --user=user                                                                      Comma separated list of users to
                                                                                       add.
```

_See code: [src/commands/chat/user/create.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/user/create.js)_

## `stream chat:user:flag`

```
USAGE
  $ stream chat:user:flag

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The ID of the offending user.
```

_See code: [src/commands/chat/user/flag.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/user/flag.js)_

## `stream chat:user:mute`

```
USAGE
  $ stream chat:user:mute

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The unique identifier of the user to mute.
```

_See code: [src/commands/chat/user/mute.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/user/mute.js)_

## `stream chat:user:remove`

```
USAGE
  $ stream chat:user:remove

OPTIONS
  -c, --channel=channel        Channel name.
  -j, --json                   Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --moderators=moderators  (required) Comma separated list of moderators to remove.
  -t, --type=type              Channel type.
```

_See code: [src/commands/chat/user/remove.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.52/src/commands/chat/user/remove.js)_
