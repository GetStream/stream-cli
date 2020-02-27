"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _uuid = require("uuid");

var _chatAuth = require("../../../utils/auth/chat-auth");

var _config = require("../../../utils/config");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class MessageCreate extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(MessageCreate);

    try {
      const {
        name
      } = await (0, _config.credentials)(this);

      if (!flags.channel) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'user',
          message: `What is the unique identifier for the user sending this message?`,
          default: (0, _uuid.v4)(),
          required: true
        }, {
          type: 'input',
          name: 'name',
          message: `What is the name of the user sending this message?`,
          default: name,
          required: true
        }, {
          type: 'input',
          name: 'image',
          message: `What is an absolute URL to the avatar of the user sending this message?`,
          hint: 'optional',
          required: false
        }, {
          type: 'input',
          name: 'channel',
          message: `What is the unique identifier for the channel?`,
          required: true
        }, {
          type: 'select',
          name: 'type',
          message: 'What type of channel is this?',
          required: true,
          choices: [{
            message: 'Livestream',
            value: 'livestream'
          }, {
            message: 'Messaging',
            value: 'messaging'
          }, {
            message: 'Gaming',
            value: 'gaming'
          }, {
            message: 'Commerce',
            value: 'commerce'
          }, {
            message: 'Team',
            value: 'team'
          }]
        }, {
          type: 'input',
          name: 'message',
          message: 'What is the message you would like to send?',
          required: true
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const channel = await client.channel(flags.type, flags.channel);
      const create = await channel.sendMessage({
        text: flags.message,
        user: {
          id: flags.user,
          name: flags.name,
          image: flags.image || null
        }
      });

      if (flags.json) {
        this.log(JSON.stringify(create.message));
        this.exit();
      }

      this.log(`Message ${_chalk.default.bold(create.message.id)} was created.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

MessageCreate.flags = {
  user: _command.flags.string({
    char: 'u',
    description: 'The ID of the user sending the message.',
    required: false
  }),
  name: _command.flags.string({
    char: 'n',
    description: 'The name of the user sending the message.',
    required: false
  }),
  image: _command.flags.string({
    char: 'i',
    description: 'Absolute URL for an avatar of the user sending the message.',
    required: false
  }),
  type: _command.flags.string({
    char: 't',
    description: 'The type of channel.',
    required: false
  }),
  channel: _command.flags.string({
    char: 'c',
    description: 'The ID of the channel that you would like to send a message to.',
    required: false
  }),
  message: _command.flags.string({
    char: 'm',
    description: 'The message you would like to send as plaintext.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
MessageCreate.description = 'Creates a new message.';
module.exports.MessageCreate = MessageCreate;