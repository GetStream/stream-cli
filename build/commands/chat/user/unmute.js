"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class UserUnmute extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(UserUnmute);

    try {
      if (!flags.user) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'user',
          message: 'What is the unique identifier for the user?',
          required: true
        }]);
        flags.user = res.user;
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const response = await client.unmuteUser(flags.user);

      if (flags.json) {
        this.log(JSON.stringify(response));
        this.exit();
      }

      this.log(`User ${_chalk.default.bold(flags.user)} has been unmuted.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

UserUnmute.flags = {
  user: _command.flags.string({
    char: 'u',
    description: 'The unique identifier of the user to unmute.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
UserUnmute.description = 'Unmutes a user.';
module.exports.UserUnmute = UserUnmute;