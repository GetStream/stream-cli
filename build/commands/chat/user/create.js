"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _uuid = require("uuid");

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class UserCreate extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(UserCreate);

    try {
      if (!flags.user || !flags.role) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'user',
          message: 'What is the unique identifier for the user?',
          default: (0, _uuid.v4)(),
          required: true
        }, {
          type: 'select',
          name: 'role',
          message: 'What role would you like assign to the user?',
          required: true,
          choices: [{
            message: 'Admin',
            value: 'admin'
          }, {
            message: 'User',
            value: 'user'
          }]
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const create = await client.updateUser({
        id: flags.user,
        role: flags.role
      });

      if (flags.json) {
        this.log(JSON.stringify(create));
        this.exit();
      }

      this.log(`The user ${_chalk.default.bold(flags.user)} (${flags.role}) has been created.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

UserCreate.flags = {
  user: _command.flags.string({
    char: 'u',
    description: 'Comma separated list of users to add.',
    required: false
  }),
  role: _command.flags.string({
    char: 'r',
    description: 'The role to assign to the user.',
    options: ['admin', 'user'],
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
UserCreate.description = 'Creates a new user.';
module.exports.UserCreate = UserCreate;