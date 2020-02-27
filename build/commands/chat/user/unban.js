"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class UserUnban extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(UserUnban);

    try {
      if (!flags.user) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'user',
          message: 'What is the unique identifier for the user?',
          required: true
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const unban = await client.unbanUser(flags.user);

      if (flags.json) {
        this.log(JSON.stringify(unban));
        this.exit();
      }

      this.log(`The user ${_chalk.default.bold(flags.user)} has been unbanned.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

UserUnban.flags = {
  user: _command.flags.string({
    char: 'u',
    description: 'The unique identifier of the user to unban.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
UserUnban.description = 'Unbans a user.';
module.exports.UserUnban = UserUnban;