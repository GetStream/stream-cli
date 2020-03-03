"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

var _config = require("../../../utils/config");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ChannelType extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(ChannelType);

    try {
      if (!flags.type) {
        const res = await (0, _enquirer.prompt)([{
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
          name: 'commands',
          message: 'What custom commands would you like to allow?',
          required: false
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const payload = {};

      for (const key in flags) {
        if (res.hasOwnProperty(key)) {
          payload[key] = flags[key];
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const update = await client.updateChannelType(flags.type, ...payload);

      if (flags.json) {
        this.log(JSON.stringify(update));
        this.exit();
      }

      this.log(`Channel ${_chalk.default.bold(flags.channel)} type has been updated.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

ChannelType.flags = {
  channel: _command.flags.string({
    char: 'c',
    description: 'The ID of the channel you wish to update.',
    required: false
  }),
  type: _command.flags.string({
    char: 't',
    description: 'Type of channel.',
    required: false
  }),
  name: _command.flags.string({
    char: 'n',
    description: 'Name of the channel room.',
    required: false
  }),
  image: _command.flags.string({
    char: 'i',
    description: 'URL to the channel image.',
    required: false
  }),
  description: _command.flags.string({
    char: 'd',
    description: 'Description for the channel.',
    required: false
  }),
  reason: _command.flags.string({
    char: 'r',
    description: 'Reason for changing channel.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
ChannelType.description = 'Updates a channels type configuration.';
module.exports.ChannelType = ChannelType;