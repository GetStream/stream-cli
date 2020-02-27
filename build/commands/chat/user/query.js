"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chatAuth = require("../../../utils/auth/chat-auth");

class UserQuery extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(UserQuery);

    try {
      if (!flags.query) {
        const query = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'query',
          message: 'What is the query you would like to perform?',
          hint: 'optional',
          required: false
        }]);
        flags.query = query;
      }

      if (!flags.sort) {
        const sort = await (0, _enquirer.prompt)([{
          type: 'select',
          name: 'sort',
          message: 'What user key would you like to sort on?',
          choices: [{
            message: 'ID',
            value: 'id'
          }, {
            message: 'Last Active',
            value: 'last_active'
          }, {
            message: 'None',
            value: 'none'
          }],
          required: false
        }]);
        flags.sort = sort;
      }

      if (!flags.limit) {
        const limit = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'query',
          message: 'How many records would you like to display?',
          default: 25,
          required: false
        }]);
        flags.limit = limit;
      }

      if (!flags.offset) {
        const offset = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'query',
          message: 'How many users would you like to skip?',
          hint: 0,
          default: 0,
          required: false
        }]);
        flags.offset = offset;
      }

      const s = {};

      if (flags.sort.sort !== 'none') {
        switch (flags.sort) {
          case 'id':
            s.id = -1;
            break;

          case 'last_active':
            s.last_active = -1;
            break;

          default:
            s.id = -1;
        }
      }

      let q = {};

      if (flags.query.query.length) {
        q = JSON.parse(flags.query.query.replace(/\'/g, ''));
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const {
        users
      } = await client.queryUsers(q, s, {
        limit: parseInt(flags.limit.query, 10) || 25,
        offset: parseInt(flags.offset.query, 10) || 0
      });

      if (flags.json) {
        this.log(JSON.stringify(users));
        this.exit();
      }

      this.log(users);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

UserQuery.flags = {
  query: _command.flags.string({
    char: 'q',
    description: 'The query you would like to perform.',
    required: false
  }),
  sort: _command.flags.string({
    char: 's',
    description: 'Display the current status of the user.',
    required: false
  }),
  limit: _command.flags.string({
    char: 'l',
    description: 'The limit to apply to the query.',
    required: false
  }),
  offset: _command.flags.string({
    char: 'o',
    description: 'The offset to apply to the query.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
UserQuery.description = 'Queries all users.';
module.exports.UserQuery = UserQuery;