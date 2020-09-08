`stream chat`
=============

Adds a member to a channel.

* [`stream chat:channel:add_member`](#stream-chatchanneladd_member)
* [`stream chat:channel:create`](#stream-chatchannelcreate)
* [`stream chat:channel:demote_moderator`](#stream-chatchanneldemote_moderator)
* [`stream chat:channel:get`](#stream-chatchannelget)
* [`stream chat:channel:hide`](#stream-chatchannelhide)
* [`stream chat:channel:list`](#stream-chatchannellist)
* [`stream chat:channel:promote_moderator`](#stream-chatchannelpromote_moderator)
* [`stream chat:channel:query`](#stream-chatchannelquery)
* [`stream chat:channel:remove`](#stream-chatchannelremove)
* [`stream chat:channel:show`](#stream-chatchannelshow)
* [`stream chat:channel:type`](#stream-chatchanneltype)
* [`stream chat:channel:update`](#stream-chatchannelupdate)
* [`stream chat:channel:update_type`](#stream-chatchannelupdate_type)
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
* [`stream chat:user:deactivate`](#stream-chatuserdeactivate)
* [`stream chat:user:flag`](#stream-chatuserflag)
* [`stream chat:user:get`](#stream-chatuserget)
* [`stream chat:user:mute`](#stream-chatusermute)
* [`stream chat:user:query`](#stream-chatuserquery)
* [`stream chat:user:reactivate`](#stream-chatuserreactivate)
* [`stream chat:user:remove`](#stream-chatuserremove)
* [`stream chat:user:unban`](#stream-chatuserunban)
* [`stream chat:user:unflag`](#stream-chatuserunflag)
* [`stream chat:user:unmute`](#stream-chatuserunmute)
* [`stream chat:user:update`](#stream-chatuserupdate)

## `stream chat:channel:add_member`

Adds a member to a channel.

```
USAGE
  $ stream chat:channel:add_member

OPTIONS
  -c, --channel=channel  A unique ID for the channel add the user to.
  -i, --image=image      URL to channel image.
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -n, --name=name        Name of the channel room.
  -r, --data=data        The role of the user you are adding.
  -t, --type=type        Type of channel.
  -u, --users=users      Unique identifier for the user you are adding.
```

## `stream chat:channel:create`

Creates a new channel.

```
USAGE
  $ stream chat:channel:create

OPTIONS
  -c, --channel=channel  [default: 0da0bfd2-8ebd-4645-8710-9cd50b5613df] A unique ID for the channel you wish to create.
  -d, --data=data        Additional data as JSON.
  -i, --image=image      URL to channel image.
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -n, --name=name        Name of the channel room.
  -t, --type=type        Type of channel.
  -u, --users=users      Comma separated list of users to add.
```

## `stream chat:channel:demote_moderator`

Demotes a moderator from a channel.

```
USAGE
  $ stream chat:channel:demote_moderator

OPTIONS
  -c, --channel=channel  A unique ID for the channel you wish to create.
  -t, --type=type        Type of channel.
  -u, --user=user        A unique ID for user to demote from a moderator.
```

## `stream chat:channel:get`

Gets a specific channel by its ID and type.

```
USAGE
  $ stream chat:channel:get

OPTIONS
  -c, --channel=channel  The channel ID you wish to retrieve.
  -t, --type=type        Type of channel.
```

## `stream chat:channel:hide`

Hides a channel.

```
USAGE
  $ stream chat:channel:hide

OPTIONS
  -c, --channel=channel  The channel ID you wish to remove.
  -t, --type=type        Type of channel.
```

## `stream chat:channel:list`

Lists all channels.

```
USAGE
  $ stream chat:channel:list

OPTIONS
  -l, --limit=limit    (required) Channel list limit.
  -o, --offset=offset  (required) Channel list offset.
```

## `stream chat:channel:promote_moderator`

Promotes a user to a moderator in a channel.

```
USAGE
  $ stream chat:channel:promote_moderator

OPTIONS
  -c, --channel=channel  A unique ID for the channel you wish to create.
  -t, --type=type        Type of channel.
  -u, --user=user        A unique ID for user user to demote.
```

## `stream chat:channel:query`

Queries all channels.

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

Removes a channel.

```
USAGE
  $ stream chat:channel:remove

OPTIONS
  -c, --channel=channel  The channel ID you wish to remove.
  -t, --type=type        Type of channel.
```

## `stream chat:channel:show`

Unhindes (shows) a channel.

```
USAGE
  $ stream chat:channel:show

OPTIONS
  -c, --channel=channel  The channel ID you wish to remove.
  -t, --type=type        Type of channel.
```

## `stream chat:channel:type`

Updates a channels type configuration.

```
USAGE
  $ stream chat:channel:type

OPTIONS
  -c, --channel=channel          The ID of the channel you wish to update.
  -d, --description=description  Description for the channel.
  -i, --image=image              URL to the channel image.
  -j, --json                     Output results in JSON. When not specified, returns output in a human friendly format.
  -n, --name=name                Name of the channel room.
  -r, --reason=reason            Reason for changing channel.
  -t, --type=type                Type of channel.
```

## `stream chat:channel:update`

Updates a channel.

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

## `stream chat:channel:update_type`

Updates a channels type configuration.

```
USAGE
  $ stream chat:channel:update_type

OPTIONS
  -a, --automod=automod                      Enable or disable automod (enabled/disabled)
  -a, --message_retention=message_retention  How long to retain messages (defaults to infinite)
  -a, --reactions                            Enable or disable reactions (true/false)
  -c, --connect_events                       Enable or disable connect events (true/false)
  -e, --read_events=read_events              Enable or disable read events (true/false)

  -j, --json                                 Output results in JSON. When not specified, returns output in a human
                                             friendly format.

  -m, --mutes                                Enable or disable mutes (true/false)

  -p, --replies                              Enable or disable replies (true/false)

  -s, --search                               Enable or disable search (true/false)

  -t, --type=type                            Type of channel.

  -y, --typing_events                        Enable or disable typing events (true/false)
```

## `stream chat:log`

Logs events in realtime.

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

Creates a new message.

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

Flags a message.

```
USAGE
  $ stream chat:message:flag

OPTIONS
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message  The unique identifier of the message you want to flag.
```

## `stream chat:message:list`

Lists all messages.

```
USAGE
  $ stream chat:message:list

OPTIONS
  -c, --channel=channel  The ID of the channel that you would like to send a message to.
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -t, --type=type        The type of channel.
```

## `stream chat:message:remove`

Removes a message.

```
USAGE
  $ stream chat:message:remove

OPTIONS
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message  The unique identifier of the message you would like to remove.
```

## `stream chat:message:unflag`

Unflags a message.

```
USAGE
  $ stream chat:message:unflag

OPTIONS
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message  The unique identifier of the message you want to flag.
```

## `stream chat:message:update`

Updates a message.

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

Specifies APN for push notifications.

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

Adds a new device for push.

```
USAGE
  $ stream chat:push:device:add

OPTIONS
  -d, --device_id=device_id  Device id or token.
  -p, --provider=provider    Push provider
  -u, --user_id=user_id      User ID
```

## `stream chat:push:device:delete`

Removes a device from push.

```
USAGE
  $ stream chat:push:device:delete

OPTIONS
  -d, --device_id=device_id  Device id or token.
  -u, --user_id=user_id      User ID
```

## `stream chat:push:device:get`

Gets all devices registered for push.

```
USAGE
  $ stream chat:push:device:get

OPTIONS
  -u, --user_id=user_id  User ID
```

## `stream chat:push:firebase`

Specifies Firebase for push notifications.

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

Gets push notification settings.

```
USAGE
  $ stream chat:push:get

OPTIONS
  -j, --json  Output results in JSON. When not specified, returns output in a human friendly format.
```

## `stream chat:push:test`

Tests push notifications.

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

Tests webhook notifications.

```
USAGE
  $ stream chat:push:webhook

OPTIONS
  -j, --json     Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --url=url  A fully qualified URL for webhook support.
```

## `stream chat:reaction:create`

Creates a new reaction.

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

Removes a reaction.

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

Bans a user.

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

Creates a new user.

```
USAGE
  $ stream chat:user:create

OPTIONS
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -r, --role=admin|user  The role to assign to the user.
  -u, --user=user        Comma separated list of users to add.
```

## `stream chat:user:deactivate`

Allows for deactivating a user and wiping all of their messages.

```
USAGE
  $ stream chat:user:deactivate

OPTIONS
  -h, --hard=hard  Hard deletes all messages associated with the user.
  -j, --json=json  Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --user=user  A unique ID of the user you would like to deactivate.
```

## `stream chat:user:flag`

Flags a user for bad behavior.

```
USAGE
  $ stream chat:user:flag

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The ID of the offending user.
```

## `stream chat:user:get`

Get a user by their unique ID.

```
USAGE
  $ stream chat:user:get

OPTIONS
  -j, --json               Output results in JSON. When not specified, returns output in a human friendly format.
  -p, --presence=presence  Display the current status of the user.
  -u, --user=user          The unique identifier of the user to get.
```

## `stream chat:user:mute`

Mutes a user.

```
USAGE
  $ stream chat:user:mute

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The unique identifier of the user to mute.
```

## `stream chat:user:query`

Queries all users.

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

## `stream chat:user:reactivate`

Reactivates a user who was previously deactivated.

```
USAGE
  $ stream chat:user:reactivate

OPTIONS
  -j, --json=json        Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --user=user        A unique ID of the user you would like to reactivate.
  -r, --restore=restore  Restores all deleted messages associated with the user.
```

## `stream chat:user:remove`

Allows for deactivating a user and wiping all of their messages.

```
USAGE
  $ stream chat:user:remove

OPTIONS
  -j, --json=json  Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --user=user  A unique ID of the user you would like to remove.
```

## `stream chat:user:unban`

Unbans a user.

```
USAGE
  $ stream chat:user:unban

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The unique identifier of the user to unban.
```

## `stream chat:user:unflag`

Unflags a user.

```
USAGE
  $ stream chat:user:unflag

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The ID of the offending user.
```

## `stream chat:user:unmute`

Unmutes a user.

```
USAGE
  $ stream chat:user:unmute

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  The unique identifier of the user to unmute.
```

## `stream chat:user:update`

Updates a user.

```
USAGE
  $ stream chat:user:update

OPTIONS
  -i, --id=id        The unique identifier for the user.
  -m, --image=image  URL to the image of the user.
  -n, --name=name    Name of the user.
```
