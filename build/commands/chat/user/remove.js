"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class UserRemove extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(UserRemove);

    try {
      if (!flags.user) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'user',
          message: 'What is the unique ID of the user you would like to remove?',
          required: true
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const exists = await client.queryUsers({
        id: flags.user
      }, {
        id: -1
      });

      if (!exists.users.length) {
        this.log(`User ${flags.user} does not exist or has already been removed.`);
        this.exit();
      }

      const remove = await client.deleteUser(flags.user, {
        mark_messages_deleted: true,
        hard_delete: true
      });

      if (flags.json) {
        this.log(JSON.stringify(remove));
        this.exit();
      }

      this.log(`${_chalk.default.bold(flags.user)} has been removed.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

UserRemove.flags = {
  user: _command.flags.string({
    char: 'm',
    description: 'A unique ID of the user you would like to remove.',
    required: false
  }),
  json: _command.flags.string({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
UserRemove.description = 'Allows for deactivating a user and wiping all of their messages.';
module.exports.UserRemove = UserRemove;