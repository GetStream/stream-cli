"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class UserReactivate extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(UserReactivate);

    try {
      if (!flags.user || !flags.restore) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'user',
          message: 'What is the unique ID of the user you would like to reactivate?',
          required: true
        }, {
          type: 'select',
          name: 'restore',
          message: 'Would you like to restore all messages?',
          required: true,
          choices: [{
            message: 'No',
            value: false
          }, {
            message: 'Yes',
            value: true
          }]
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const {
        user
      } = await client.reactivateUser(flags.user, {
        restore_messages: Boolean(flags.restore)
      });

      if (flags.json) {
        this.log(JSON.stringify(user));
        this.exit();
      }

      this.log(`${_chalk.default.bold(flags.user)} has been reactivated.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

UserReactivate.flags = {
  user: _command.flags.string({
    char: 'm',
    description: 'A unique ID of the user you would like to reactivate.',
    required: false
  }),
  restore: _command.flags.string({
    char: 'r',
    description: 'Restores all deleted messages associated with the user.',
    required: false
  }),
  json: _command.flags.string({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
UserReactivate.description = 'Reactivates a user who was previously deactivated.';
module.exports.UserReactivate = UserReactivate;