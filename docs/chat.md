`stream chat`
=============

configure and manage all things related to chat

* [`stream chat:channel:edit`](#stream-chatchanneledit)
* [`stream chat:channel:get`](#stream-chatchannelget)
* [`stream chat:channel:init`](#stream-chatchannelinit)
* [`stream chat:channel:list`](#stream-chatchannellist)
* [`stream chat:channel:query`](#stream-chatchannelquery)
* [`stream chat:log`](#stream-chatlog)
* [`stream chat:message:remove`](#stream-chatmessageremove)
* [`stream chat:message:send`](#stream-chatmessagesend)
* [`stream chat:moderate:ban`](#stream-chatmoderateban)
* [`stream chat:moderate:flag`](#stream-chatmoderateflag)
* [`stream chat:moderate:mute`](#stream-chatmoderatemute)
* [`stream chat:user:add`](#stream-chatuseradd)
* [`stream chat:user:remove`](#stream-chatuserremove)

## `stream chat:channel:edit`

```
USAGE
  $ stream chat:channel:edit

OPTIONS
  -d, --data=data                                       Additional data as JSON.
  -i, --id=id                                           (required) The ID of the channel you wish to edit.
  -m, --members=members                                 Comma separated list of members.
  -n, --name=name                                       (required) Name of the channel room.
  -r, --reason=reason                                   (required) Reason for changing channel.
  -t, --type=livestream|messaging|gaming|commerce|team  (required) Type of channel.
  -u, --url=url                                         URL to the channel image.
```

_See code: [src/commands/chat/channel/edit.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/channel/edit.js)_

## `stream chat:channel:get`

```
USAGE
  $ stream chat:channel:get

OPTIONS
  -i, --id=id                                           (required) The channel ID you wish to get.
  -t, --type=livestream|messaging|gaming|commerce|team  (required) Type of channel.
```

_See code: [src/commands/chat/channel/get.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/channel/get.js)_

## `stream chat:channel:init`

```
USAGE
  $ stream chat:channel:init

OPTIONS
  -d, --data=data                                       Additional data as a JSON.

  -i, --id=id                                           (required) [default: 941e3800-12ad-49be-8952-3ca7f8346727] A
                                                        unique ID for the channel you wish to create.

  -m, --members=members                                 Comma separated list of members to add to the channel.

  -n, --name=name                                       (required) Name of the channel room.

  -t, --type=livestream|messaging|gaming|commerce|team  (required) Type of channel.

  -u, --image=image                                     URL to channel image.
```

_See code: [src/commands/chat/channel/init.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/channel/init.js)_

## `stream chat:channel:list`

```
USAGE
  $ stream chat:channel:list
```

_See code: [src/commands/chat/channel/list.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/channel/list.js)_

## `stream chat:channel:query`

```
USAGE
  $ stream chat:channel:query

OPTIONS
  -f, --filter=filter                                   Filters to apply to the query.
  -i, --id=id                                           (required) The channel ID you wish to query.
  -s, --sort=sort                                       Sort to apply to the query.
  -t, --type=livestream|messaging|gaming|commerce|team  Type of channel.
```

_See code: [src/commands/chat/channel/query.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/channel/query.js)_

## `stream chat:log`

```
USAGE
  $ stream chat:log

OPTIONS
  -e, 
  --event=all|user.status.changed|user.watching.start|user.watching.stop|user.updated|typing.start|typing.stop|message.n
  ew|message.updated|message.deleted|message.seen|message.reaction|member.added|member.removed|channel.updated|health.ch
  eck|connection.changed|connection.recovered
      The type of event you want to listen on.

  -i, --id=id
      [default: 4d5c8d59-585a-4aca-9532-5c12c359107d] The channel ID you wish to log.

  -t, --type=livestream|messaging|gaming|commerce|team
      (required) The type of channel.
```

_See code: [src/commands/chat/log/index.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/log/index.js)_

## `stream chat:message:remove`

```
USAGE
  $ stream chat:message:remove

OPTIONS
  -i, --id=id  (required) The channel ID that you would like to remove.
```

_See code: [src/commands/chat/message/remove.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/message/remove.js)_

## `stream chat:message:send`

```
USAGE
  $ stream chat:message:send

OPTIONS
  -a, --attachments=attachments                         A JSON payload of attachments to send with the message.
  -i, --id=id                                           The channel ID that you would like to send a message to.
  -m, --message=message                                 (required) The message you would like to send as plaintext.
  -t, --type=livestream|messaging|gaming|commerce|team  (required) The type of channel.

  -u, --user=user                                       (required) [default: *] The ID of the acting user sending the
                                                        message.
```

_See code: [src/commands/chat/message/send.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/message/send.js)_

## `stream chat:moderate:ban`

```
USAGE
  $ stream chat:moderate:ban

OPTIONS
  -r, --reason=reason    (required) A reason for adding a timeout.
  -t, --timeout=timeout  (required) [default: 60] Duration of timeout in minutes.
  -u, --user=user        (required) The ID of the offending user.
```

_See code: [src/commands/chat/moderate/ban.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/moderate/ban.js)_

## `stream chat:moderate:flag`

```
USAGE
  $ stream chat:moderate:flag

OPTIONS
  -m, --message=message  The ID of the message you want to flag.
  -u, --user=user        The ID of the offending user.
```

_See code: [src/commands/chat/moderate/flag.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/moderate/flag.js)_

## `stream chat:moderate:mute`

```
USAGE
  $ stream chat:moderate:mute

OPTIONS
  -u, --user=user  (required) The ID of the offending user.
```

_See code: [src/commands/chat/moderate/mute.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/moderate/mute.js)_

## `stream chat:user:add`

```
USAGE
  $ stream chat:user:add

OPTIONS
  -i, --id=id                  (required) Channel name.
  -m, --moderators=moderators  (required) Comma separated list of moderators to add.
  -t, --type=type              (required) Channel type.
```

_See code: [src/commands/chat/user/add.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/user/add.js)_

## `stream chat:user:remove`

```
USAGE
  $ stream chat:user:remove

OPTIONS
  -i, --id=id                  (required) Channel name.
  -m, --moderators=moderators  (required) Comma separated list of moderators to remove.
  -t, --type=type              (required) Channel type.
```

_See code: [src/commands/chat/user/remove.js](https://github.com/getstream/stream-cli/blob/v0.0.1-beta.17/src/commands/chat/user/remove.js)_
