"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

var _config = require("../../../utils/config");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class MessageUpdate extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(MessageUpdate);

    try {
      const {
        name
      } = await (0, _config.credentials)(this);

      if (!flags.message || !flags.text) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'message',
          message: `What is the unique identifier for the message?`,
          required: true
        }, {
          type: 'input',
          name: 'text',
          message: 'What is the updated message?',
          required: true
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const payload = {
        id: flags.message,
        text: flags.text,
        user: {
          id: 'CLI',
          name
        }
      };

      if (flags.attachments) {
        payload.attachments = JSON.parse(flags.attachments);
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      await client.setUser({
        id: 'CLI',
        status: 'invisible'
      });
      const update = await client.updateMessage(payload);

      if (flags.json) {
        this.log(JSON.stringify(update));
        this.exit();
      }

      this.log(`Message ${_chalk.default.bold(flags.message)} has been updated.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

MessageUpdate.flags = {
  message: _command.flags.string({
    char: 'm',
    description: 'The unique identifier for the message.',
    required: false
  }),
  text: _command.flags.string({
    char: 't',
    description: 'The message you would like to send as text.',
    required: false
  }),
  attachments: _command.flags.string({
    char: 'a',
    description: 'A JSON payload of attachments to send along with a message.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
MessageUpdate.description = 'Updates a message.';
module.exports.MessageUpdate = MessageUpdate;