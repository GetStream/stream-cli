"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class UserGet extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(UserGet);

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
      const user = await client.queryUsers({
        id: {
          $in: [flags.user]
        }
      }, {
        id: -1
      });

      if (!user.users.length) {
        this.log(`User ${_chalk.default.bold(flags.user)} could not be found.`);
        this.exit();
      }

      if (flags.json) {
        this.log(JSON.stringify(user.users[0]));
        this.exit();
      }

      this.log(user.users[0]);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

UserGet.flags = {
  user: _command.flags.string({
    char: 'u',
    description: 'The unique identifier of the user to get.',
    required: false
  }),
  presence: _command.flags.string({
    char: 'p',
    description: 'Display the current status of the user.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
UserGet.description = 'Get a user by their unique ID.';
module.exports.UserGet = UserGet;