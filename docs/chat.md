`stream chat`
=============

configure and manage all things related to chat

* [`stream chat:channel:get`](#stream-chatchannelget)
* [`stream chat:channel:init`](#stream-chatchannelinit)
* [`stream chat:channel:list`](#stream-chatchannellist)
* [`stream chat:channel:query`](#stream-chatchannelquery)
* [`stream chat:channel:update`](#stream-chatchannelupdate)
* [`stream chat:log`](#stream-chatlog)
* [`stream chat:message:create`](#stream-chatmessagecreate)
* [`stream chat:message:remove`](#stream-chatmessageremove)
* [`stream chat:moderate:ban`](#stream-chatmoderateban)
* [`stream chat:moderate:flag`](#stream-chatmoderateflag)
* [`stream chat:moderate:mute`](#stream-chatmoderatemute)
* [`stream chat:reaction:add`](#stream-chatreactionadd)
* [`stream chat:reaction:remove`](#stream-chatreactionremove)
* [`stream chat:settings:get`](#stream-chatsettingsget)
* [`stream chat:settings:set`](#stream-chatsettingsset)
* [`stream chat:user:create`](#stream-chatusercreate)
* [`stream chat:user:remove`](#stream-chatuserremove)

## `stream chat:channel:get`

```
USAGE
  $ stream chat:channel:get

OPTIONS
  -c, --channel=channel                                 The channel ID you wish to retrieve.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -t, --type=livestream|messaging|gaming|commerce|team  Type of channel.
```

_See code: [src/commands/chat/channel/get.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/channel/get.js)_

## `stream chat:channel:init`

```
USAGE
  $ stream chat:channel:init

OPTIONS
  -c, --channel=channel                                 (required) [default: 5c74e59e-ee57-4e50-b387-f6f3816211ab] A
                                                        unique ID for the channel you wish to create.

  -d, --data=data                                       Additional data as a JSON.

  -i, --image=image                                     URL to channel image.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -m, --members=members                                 Comma separated list of members to add to the channel.

  -n, --name=name                                       (required) Name of the channel room.

  -t, --type=livestream|messaging|gaming|commerce|team  (required) Type of channel.
```

_See code: [src/commands/chat/channel/init.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/channel/init.js)_

## `stream chat:channel:list`

```
USAGE
  $ stream chat:channel:list

OPTIONS
  -j, --json  Output results in JSON. When not specified, returns output in a human friendly format.
```

_See code: [src/commands/chat/channel/list.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/channel/list.js)_

## `stream chat:channel:query`

```
USAGE
  $ stream chat:channel:query

OPTIONS
  -c, --channel=channel                                 (required) The channel ID you wish to query.
  -f, --filter=filter                                   Filters to apply to the query.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -s, --sort=sort                                       Sort to apply to the query.

  -t, --type=livestream|messaging|gaming|commerce|team  Type of channel.
```

_See code: [src/commands/chat/channel/query.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/channel/query.js)_

## `stream chat:channel:update`

```
USAGE
  $ stream chat:channel:update

OPTIONS
  -d, --data=data                                       Additional data as JSON.
  -i, --id=id                                           (required) The ID of the channel you wish to update.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -m, --members=members                                 Comma separated list of members.

  -n, --name=name                                       (required) Name of the channel room.

  -r, --reason=reason                                   (required) Reason for changing channel.

  -t, --type=livestream|messaging|gaming|commerce|team  (required) Type of channel.

  -u, --url=url                                         URL to the channel image.
```

_See code: [src/commands/chat/channel/update.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/channel/update.js)_

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
      (required) The type of channel.
```

_See code: [src/commands/chat/log/index.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/log/index.js)_

## `stream chat:message:create`

```
USAGE
  $ stream chat:message:create

OPTIONS
  -a, --attachments=attachments                         A JSON payload of attachments to send along with a message.
  -c, --channel=channel                                 The ID of the channel that you would like to send a message to.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -m, --message=message                                 The message you would like to send as plaintext.

  -t, --type=livestream|messaging|gaming|commerce|team  The type of channel.

  -u, --user=user                                       The ID of the user sending the message.
```

_See code: [src/commands/chat/message/create.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/message/create.js)_

## `stream chat:message:remove`

```
USAGE
  $ stream chat:message:remove

OPTIONS
  -c, --channel=channel  (required) The channel ID you are targeting.
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message  (required) The ID of the message you would like to remove.
```

_See code: [src/commands/chat/message/remove.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/message/remove.js)_

## `stream chat:moderate:ban`

```
USAGE
  $ stream chat:moderate:ban

OPTIONS
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -r, --reason=reason    (required) A reason for adding a timeout.
  -t, --timeout=timeout  (required) [default: 60] Duration of timeout in minutes.
  -u, --user=user        (required) The ID of the offending user.
```

_See code: [src/commands/chat/moderate/ban.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/moderate/ban.js)_

## `stream chat:moderate:flag`

```
USAGE
  $ stream chat:moderate:flag

OPTIONS
  -j, --json             Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --message=message  The ID of the message you want to flag.
  -u, --user=user        The ID of the offending user.
```

_See code: [src/commands/chat/moderate/flag.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/moderate/flag.js)_

## `stream chat:moderate:mute`

```
USAGE
  $ stream chat:moderate:mute

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -u, --user=user  (required) The ID of the user to mute.
```

_See code: [src/commands/chat/moderate/mute.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/moderate/mute.js)_

## `stream chat:reaction:add`

```
USAGE
  $ stream chat:reaction:add

OPTIONS
  -c, --channel=channel                                 The unique identifier for the channel.
  -c, --message=message                                 The unique identifier for the message.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -r, --reaction=reaction                               A reaction for the message (e.g. love).

  -t, --type=livestream|messaging|gaming|commerce|team  The type of channel.
```

_See code: [src/commands/chat/reaction/add.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/reaction/add.js)_

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

_See code: [src/commands/chat/reaction/remove.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/reaction/remove.js)_

## `stream chat:settings:get`

```
USAGE
  $ stream chat:settings:get

OPTIONS
  -j, --json  Output results in JSON. When not specified, returns output in a human friendly format.
```

_See code: [src/commands/chat/settings/get.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/settings/get.js)_

## `stream chat:settings:set`

```
USAGE
  $ stream chat:settings:set

OPTIONS
  -j, --json       Output results in JSON. When not specified, returns output in a human friendly format.
  -n, --name=name  Full name for settings.
  -p, --p12=p12    A .p12 file for push notifications.
  -p, --pem=pem    The .pem file associated with your .p12
```

_See code: [src/commands/chat/settings/set.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/settings/set.js)_

## `stream chat:user:create`

```
USAGE
  $ stream chat:user:create

OPTIONS
  -c, --channel=channel                                 Channel identifier.

  -j, --json                                            Output results in JSON. When not specified, returns output in a
                                                        human friendly format.

  -m, --moderators=moderators                           Comma separated list of moderators.

  -t, --type=livestream|messaging|gaming|commerce|team  The type of channel.
```

_See code: [src/commands/chat/user/create.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/user/create.js)_

## `stream chat:user:remove`

```
USAGE
  $ stream chat:user:remove

OPTIONS
  -c, --channel=channel        (required) Channel name.
  -j, --json                   Output results in JSON. When not specified, returns output in a human friendly format.
  -m, --moderators=moderators  (required) Comma separated list of moderators to remove.
  -t, --type=type              (required) Channel type.
```

_See code: [src/commands/chat/user/remove.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.24/src/commands/chat/user/remove.js)_
