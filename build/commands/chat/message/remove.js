"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class MessageRemove extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(MessageRemove);

    try {
      if (!flags.message) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'message',
          message: `What is the unique identifier for the message?`,
          required: true
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const remove = await client.deleteMessage(flags.message);

      if (flags.json) {
        this.log(JSON.stringify(remove));
        this.exit();
      }

      this.log(`The message ${_chalk.default.bold(flags.message)} has been removed.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

MessageRemove.flags = {
  message: _command.flags.string({
    char: 'message',
    description: 'The unique identifier of the message you would like to remove.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
MessageRemove.description = 'Removes a message.';
module.exports.MessageRemove = MessageRemove;