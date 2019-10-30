`stream chat`
=============

configure and manage all things related to chat

* [`stream chat:channel:create`](#stream-chatchannelcreate)
* [`stream chat:channel:get`](#stream-chatchannelget)
* [`stream chat:channel:list`](#stream-chatchannellist)
* [`stream chat:channel:query`](#stream-chatchannelquery)
* [`stream chat:channel:remove`](#stream-chatchannelremove)
* [`stream chat:channel:update`](#stream-chatchannelupdate)
* [`stream chat:log`](#stream-chatlog)
* [`stream chat:message:create`](#stream-chatmessagecreate)
* [`stream chat:message:flag`](#stream-chatmessageflag)
* [`stream chat:message:list`](#stream-chatmessagelist)
* [`stream chat:message:remove`](#stream-chatmessageremove)
* [`stream chat:message:unflag`](#stream-chatmessageunflag)
* [`stream chat:message:update`](#stream-chatmessageupdate)
* [`stream chat:push:apn`](#stream-chatpushapn)
* [`stream chat:push:device:add`](#stream-chatpushdeviceadd)
* [`stream chat:push:device:delete`](#stream-chatpushdevicedelete)
* [`stream chat:push:device:get`](#stream-chatpushdeviceget)
* [`stream chat:push:firebase`](#stream-chatpushfirebase)
* [`stream chat:push:get`](#stream-chatpushget)
* [`stream chat:push:test`](#stream-chatpushtest)
* [`stream chat:push:webhook`](#stream-chatpushwebhook)
* [`stream chat:reaction:create`](#stream-chatreactioncreate)
* [`stream chat:reaction:remove`](#stream-chatreactionremove)
* [`stream chat:user:ban`](#stream-chatuserban)
* [`stream chat:user:create`](#stream-chatusercreate)
* [`stream chat:user:flag`](#stream-chatuserflag)
* [`stream chat:user:get`](#stream-chatuserget)
* [`stream chat:user:mute`](#stream-chatusermute)
* [`stream chat:user:query`](#stream-chatuserquery)
* [`stream chat:user:remove`](#stream-chatuserremove)
* [`stream chat:user:set`](#stream-chatuserset)
* [`stream chat:user:unban`](#stream-chatuserunban)
* [`stream chat:user:unflag`](#stream-chatuserunflag)
* [`stream chat:user:unmute`](#stream-chatuserunmute)

## `stream chat:channel:create`

```
USAGE
  $ stream chat:channel:create

OPTIONS
  -c, --channel=channel  [default: 032aca2a-2d0a-4461-b55b-5419872133e3] A unique ID for the channel you wish to create.
  -d, --data=data        Additional data as JSON.
  -i, --image=image      URL to channel image.
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -n, --name=name        Name of the channel room.
  -t, --type=type        Type of channel.
  -u, --users=users      Comma separated list of users to add.
```

## `stream chat:channel:get`

```
USAGE
  $ stream chat:channel:get

OPTIONS
  -c, --channel=channel  The channel ID you wish to retrieve.
  -t, --type=type        Type of channel.
```

## `stream chat:channel:list`

```
USAGE
  $ stream chat:channel:list

OPTIONS
  -l, --limit=limit    (required) Channel list limit.
  -o, --offset=offset  (required) Channel list offset.
```

## `stream chat:channel:query`

```
USAGE
  $ stream chat:channel:query

OPTIONS
  -c, --channel=channel  The unique identifier for the channel you want to query.
  -f, --filter=filter    Filters to apply to the query.
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -s, --sort=sort        Sort to apply to the query.
  -t, --type=type        Type of channel.
```

## `stream chat:channel:remove`

```
USAGE
  $ stream chat:channel:remove

OPTIONS
  -c, --channel=channel  The channel ID you wish to remove.
  -t, --type=type        Type of channel.
```

## `stream chat:channel:update`

```
USAGE
  $ stream chat:channel:update

OPTIONS
  -c, --channel=channel          The ID of the channel you wish to update.
  -d, --description=description  Description for the channel.
  -i, --image=image              URL to the channel image.
  -j, --json                     Output results in JSON. When not specified, returns output in a human friendly format.
  -n, --name=name                Name of the channel room.
  -r, --reason=reason            Reason for changing channel.
  -t, --type=type                Type of channel.
```

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

  -t, --type=type
      The type of channel.
```

## `stream chat:message:create`

```
USAGE
  $ stream chat:message:create

OPTIONS
  -c, --channel=channel  The ID of the channel that you would like to send a message to.
  -i, --image=image      Absolute URL for an avatar of the user sending the message.
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message  The message you would like to send as plaintext.
  -n, --name=name        The name of the user sending the message.
  -t, --type=type        The type of channel.
  -u, --user=user        The ID of the user sending the message.
```

## `stream chat:message:flag`

```
USAGE
  $ stream chat:message:flag

OPTIONS
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message  The unique identifier of the message you want to flag.
```

## `stream chat:message:list`

```
USAGE
  $ stream chat:message:list

OPTIONS
  -c, --channel=channel  The ID of the channel that you would like to send a message to.
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -t, --type=type        The type of channel.
```

## `stream chat:message:remove`

```
USAGE
  $ stream chat:message:remove

OPTIONS
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message  The unique identifier of the message you would like to remove.
```

## `stream chat:message:unflag`

```
USAGE
  $ stream chat:message:unflag

OPTIONS
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message  The unique identifier of the message you want to flag.
```

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

## `stream chat:push:device:add`

```
USAGE
  $ stream chat:push:device:add

OPTIONS
  -d, --device_id=device_id  Device id or token.
  -p, --provider=provider    Push provider
  -u, --user_id=user_id      User ID
```

## `stream chat:push:device:delete`

```
USAGE
  $ stream chat:push:device:delete

OPTIONS
  -d, --device_id=device_id  Device id or token.
  -u, --user_id=user_id      User ID
```

## `stream chat:push:device:get`

```
USAGE
  $ stream chat:push:device:get

OPTIONS
  -u, --user_id=user_id  User ID
```

## `stream chat:push:firebase`

```
USAGE
  $ stream chat:push:firebase

OPTIONS
  -d, --data_template=data_template                  JSON data template.

  -j, --json                                         Output results in JSON. When not specified, returns output in a
                                                     human friendly format.

  -k, --key=key                                      Server key for Firebase.

  -n, --notification_template=notification_template  JSON notification template.

  --disable                                          Disable Firebase push notifications and clear config.
```

## `stream chat:push:get`

```
USAGE
  $ stream chat:push:get

OPTIONS
  -j, --json  Output results in JSON. When not specified, returns output in a human friendly format.
```

## `stream chat:push:test`

```
USAGE
  $ stream chat:push:test

OPTIONS
  -a, --apn_notification_template=apn_notification_template            APN notification template
  -d, --firebase_data_template=firebase_data_template                  Firebase data template
  -f, --firebase_notification_template=firebase_notification_template  Firebase notification template

  -j, --json                                                           Output results in JSON. When not specified,
                                                                       returns output in a human friendly format.

  -m, --message_id=message_id                                          Message ID.

  -u, --user_id=user_id                                                User ID
```

## `stream chat:push:webhook`

```
USAGE
  $ stream chat:push:webhook

OPTIONS
  -j, --json     Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --url=url  A fully qualified URL for webhook support.
```

## `stream chat:reaction:create`

```
USAGE
  $ stream chat:reaction:create

OPTIONS
  -c, --channel=channel    The unique identifier for the channel.
  -c, --message=message    The unique identifier for the message.
  -j, --json               Output results in JSON. When not specified, returns output in a human friendly format.
  -r, --reaction=reaction  A reaction for the message (e.g. love).
  -t, --type=type          The type of channel.
```

## `stream chat:reaction:remove`

```
USAGE
  $ stream chat:reaction:remove

OPTIONS
  -c, --channel=channel    The unique identifier for the channel.
  -c, --message=message    The unique identifier for the message.
  -j, --json               Output results in JSON. When not specified, returns output in a human friendly format.
  -r, --reaction=reaction  The unique identifier for the reaction.
  -t, --type=type          The type of channel.
```

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

## `stream chat:user:create`

```
USAGE
  $ stream chat:user:create

OPTIONS
  -c, --channel=channel                                              Channel identifier.

  -j, --json                                                         Output results in JSON. When not specified, returns
                                                                     output in a human friendly format.

  -r, --role=admin|guest|channel_member|channel_owner|message_owner  The role to assign to the user.

  -t, --type=type                                                    The type of channel.

  -u, --user=user                                                    Comma separated list of users to add.
```

## `stream chat:user:flag`

```
USAGE
  $ stream chat:user:flag

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The ID of the offending user.
```

## `stream chat:user:get`

```
USAGE
  $ stream chat:user:get

OPTIONS
  -j, --json               Output results in JSON. When not specified, returns output in a human friendly format.
  -p, --presence=presence  Display the current status of the user.
  -u, --user=user          The unique identifier of the user to get.
```

## `stream chat:user:mute`

```
USAGE
  $ stream chat:user:mute

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The unique identifier of the user to mute.
```

## `stream chat:user:query`

```
USAGE
  $ stream chat:user:query

OPTIONS
  -j, --json           Output results in JSON. When not specified, returns output in a human friendly format.
  -l, --limit=limit    The limit to apply to the query.
  -o, --offset=offset  The offset to apply to the query.
  -q, --query=query    The query you would like to perform.
  -s, --sort=sort      Display the current status of the user.
```

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

## `stream chat:user:set`

```
USAGE
  $ stream chat:user:set

OPTIONS
  -i, --id=id        The unique identifier for the user.
  -m, --image=image  URL to the image of the user.
  -n, --name=name    Name of the user.
```

## `stream chat:user:unban`

```
USAGE
  $ stream chat:user:unban

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The unique identifier of the user to unban.
```

## `stream chat:user:unflag`

```
USAGE
  $ stream chat:user:unflag

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The ID of the offending user.
```

## `stream chat:user:unmute`

```
USAGE
  $ stream chat:user:unmute

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The unique identifier of the user to unmute.
```
