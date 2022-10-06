## stream-cli

Stream CLI Imports

### Introduction

Stream CLI allows you to manage your Stream Chat Imports easily and validate your imports quickly before uploading them
for processing.

### Validation

To import data into your Stream app, you'll need to provide a valid import file. The file needs to be structured as a
list of JSON objects, representing each item to be imported.

<details markdown="1">
<summary>Here's an example of a valid import file</summary>

```json
[
  {
    "type": "user",
    "item": {
      "id": "hannibal",
      "role": "user",
      "teams": [
        "A-team"
      ],
      "channel_mutes": [
        "messaging:HQ"
      ]
    }
  },
  {
    "type": "user",
    "item": {
      "id": "murdock",
      "role": "user",
      "teams": [
        "A-team"
      ],
      "push_notifications": {
        "disabled": true,
        "disabled_reason": "doesn't want to be disturbed"
      },
      "user_mutes": [
        "hannibal"
      ]
    }
  },
  {
    "type": "channel",
    "item": {
      "id": "HQ",
      "type": "messaging",
      "team": "A-team",
      "created_by": "hannibal"
    }
  },
  {
    "type": "member",
    "item": {
      "channel_id": "HQ",
      "channel_type": "messaging",
      "user_id": "hannibal"
    }
  },
  {
    "type": "member",
    "item": {
      "channel_id": "HQ",
      "channel_type": "messaging",
      "user_id": "murdock"
    }
  },
  {
    "type": "message",
    "item": {
      "id": "message1",
      "channel_id": "HQ",
      "channel_type": "messaging",
      "user": "hannibal",
      "text": "I love it when a plan comes together"
    }
  },
  {
    "type": "reaction",
    "item": {
      "message_id": "message1",
      "type": "like",
      "user_id": "murdock",
      "created_at": "2022-01-01T01:01:01Z"
    }
  },
  {
    "type": "device",
    "item": {
      "id": "deviceID",
      "user_id": "hannibal",
      "created_at": "2022-01-01T01:01:01Z",
      "push_provider_type": "firebase"
    }
  }
]
```
</details>

Before processing this file, we'll need to ensure it is valid. We can do this by using the validate-import command.

<details markdown="1">
<summary>`validate-import` example</summary>

```shell
$ stream-cli chat validate-import my-data.json                                                                                                                                                                                                                           9:33:17
{
  "Stats": {
    "channels": 1,
    "members": 2,
    "messages": 1,
    "reactions": 1,
    "users": 2
  },
  "Errors": null
}
```
</details>

The lack of Errors here tells us this file is valid. Additionally, the output tells us the number of each item type.
However, there are several reasons why an import file might be invalid. They can be divided into validation errors and
reference errors.

#### Validation Errors

Validation errors occur when items contain invalid data or are missing required fields.

| Error                                                                                                                                        | Reason                                                                                                         |
| -------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| `validation error: user.id required`                                                                                                         | A required field is missing                                                                                    |
| `validation error: either channel.id or channel.member_ids required`                                                                         | A channel needs either an id or a list of member ids (but not both)                                            |
| `validation error: either message.channel_id or message.channel_member_ids required`                                                         | A channel reference needs to be either a channel id or a list of channel member ids (but not both)             |
| `validation error: user.id max length exceeded (255)`                                                                                        | A value is exceeding the maximum allowed length                                                                |
| `validation error: user.id invalid ("^[@\w-]*$" allowed)`                                                                                    | A value is invalid according to a regular expression                                                           |
| `validation error: user.online is a reserved field`                                                                                          | A field is reserved and not allowed to be present in the import file                                           |
| `validation error: message.type invalid ("regular", "deleted", "system" and "reply" allowed)`                                                | A field can only be one of these values                                                                        |
| `validation error: message.type "deleted" while message.deleted_at is null`                                                                  | A message of type `deleted` must have `deleted_at` set                                                         |
| `validation error: distinct channel: ["userA" "userB"] is missing members: ["userB"]. Please include all members as separate member entries` | A distinct channel has been defined with a member `userB`, but has not been included as a separate member item |
| `validation error: duplicate user found "userA"`                                                                                             | An item occurs more than once in the import file                                                               |

#### Reference Errors

Every user, channel, or reference message needs to be included as a separate item in the import file. Additionally,
referenced user roles and channel types have to exist before importing.

| Error                                                                                                                   | Reason                                                                                          |
| ----------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `reference error: user.role "admin" doesn't exist (user "userA")`                                                       | User role must exists before importing                                                          |
| `reference error: channel.type "livestream" doesn't exist (channel "livestream:chanelA")`                               | Channel type must exist before importing                                                        |
| `reference error: channel.created_by "john" doesn't exist (channel "messaging:channelA")`                               | User `john` must be included in the import file                                                 |
| `reference error: distinct channel with type "messaging" and members:["userA" "userB"] doesn't exist`                   | The distinct channel referenced with member userA and userB must be included in the import file |
| `reference error: user "userA" with teams ["teamA"] cannot be a member of channel messaging:channelB with team "teamB"` | The user is not a part of the team assigned to this channel                                     |
| `reference error: user "userC" specified as channel member but not present in channel_members_ids: [userA userB]`       | The member items must match the members in the distinct channel definition                      |

### Uploading

Once the import file is valid, the file can be uploaded to be scheduled for import:

<details markdown="1">
<summary>`upload-import` example</summary>

```shell
$ stream-cli chat upload-import my-data.json                                                                                                                                                                                                                            10:14:36
{
  "created_at": "2022-05-16T09:02:37.991181Z",
  "path": "s3://stream-import/1171432/7e7fbaf4-e266-4877-96da-fbacf650d0a1/my-data.json",
  "mode": "upsert",
  "history": [],
  "id": "79502357-3f4b-486e-9a78-400a184a1088",
  "state": "uploaded",
  "updated_at": "2022-05-16T09:02:37.991181Z",
  "result": null,
  "size": 1230
}
```
</details>

#### Import modes

By default, the `mode` is set to `upsert`. This means that every item will be either inserted as is, or updated to match
the item as it appears in the import file. Uploads can be created with `mode` set to `insert`, to insert only and ignore
pre-existing items.

<details markdown="1">
<summary>import modes example</summary>

```shell
$ stream-cli chat upload-import valid-data.json --mode insert                                                                                                                                                                                                              11:02:44
{
  "created_at": "2022-05-16T09:06:26.356475Z",
  "path": "s3://stream-import/1171432/f9460261-ce3c-4594-a236-33bbedfa85a7/valid-data.json",
  "mode": "insert",
  "history": [],
  "id": "f0077dab-84f3-48f1-9292-2bf1b48fd6f0",
  "state": "uploaded",
  "updated_at": "2022-05-16T09:06:26.356475Z",
  "result": null,
  "size": 1230
}
```
</details>

Using import modes can be helpful when migrating to Stream Chat using a dual-write approach. 
One could do an initial import using `mode` set to `upsert`, which imports everything as is. 
Then using `mode` set to `insert`, to do additional imports to only import missing data to fill in the gaps.

### Managing imports

Once an import has been created, you can monitor its status using the `get-import` and `list-imports` commands.

`get-import` accepts an option `--watch` flag, which will periodically poll the import status.

<details markdown="1">
<summary>`get-import` example</summary>

```shell
$ stream-cli chat get-import f0077dab-84f3-48f1-9292-2bf1b48fd6f0  --watch                                                                                                                                                                                                 13:50:09
{
  "import_task": {
    "created_at": "2022-05-16T09:06:26.356475Z",
    "path": "s3://stream-import/1171432/f9460261-ce3c-4594-a236-33bbedfa85a7/valid-data.json",
    "mode": "insert",
    "history": [],
    "id": "f0077dab-84f3-48f1-9292-2bf1b48fd6f0",
    "state": "uploaded",
    "updated_at": "2022-05-16T09:06:26.356475Z",
    "result": null,
    "size": 1230
  },
  "ratelimit": {
    "limit": 300,
    "remaining": 299,
    "reset": 1652701860
  }
}
```
</details>

`list-imports` will return an imports list and accepts a `limit` and `offset` parameter.

<details markdown="1">
<summary>`list-imports` example</summary>

```shell
$ stream-cli chat list-imports                                                                                                                                                                                                                                       130 ↵ 13:51:07
[
  {
    "created_at": "2022-05-16T09:06:26.356475Z",
    "path": "s3://stream-import/1171432/f9460261-ce3c-4594-a236-33bbedfa85a7/valid-data.json",
    "mode": "insert",
    "history": [],
    "id": "f0077dab-84f3-48f1-9292-2bf1b48fd6f0",
    "state": "uploaded",
    "updated_at": "2022-05-16T09:06:26.356475Z",
    "result": null,
    "size": 1230
  },
  {
    "created_at": "2022-05-16T09:02:37.991181Z",
    "path": "s3://stream-import/1171432/7e7fbaf4-e266-4877-96da-fbacf650d0a1/valid-data.json",
    "mode": "upsert",
    "history": [],
    "id": "79502357-3f4b-486e-9a78-400a184a1088",
    "state": "uploaded",
    "updated_at": "2022-05-16T09:02:37.991181Z",
    "result": null,
    "size": 1230
  }
]
```
</details>
