"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ChannelShow extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(ChannelShow);

    try {
      if (!flags.channel || !flags.type || !flags.user) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'user',
          hint: 'user-123',
          message: 'What is the ID of the user you wish to show the channel to?',
          required: true
        }, {
          type: 'input',
          name: 'channel',
          message: `What is the unique identifier for the channel?`,
          required: true
        }, {
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
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const channel = client.channel(flags.type, flags.channel);
      await channel.show(flags.user);
      this.log(`The channel ${_chalk.default.bold(flags.channel)} has been shown.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

ChannelShow.flags = {
  user: _command.flags.string({
    char: 'u',
    description: 'User ID',
    required: false
  }),
  channel: _command.flags.string({
    char: 'c',
    description: 'The channel ID you wish to remove.',
    required: false
  }),
  type: _command.flags.string({
    char: 't',
    description: 'Type of channel.',
    required: false
  })
};
ChannelShow.description = 'Unhindes (shows) a channel.';
module.exports.ChannelShow = ChannelShow;