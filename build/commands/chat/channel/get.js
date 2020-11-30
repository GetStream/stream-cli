"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ChannelGet extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(ChannelGet);

    try {
      if (!flags.channel || !flags.type) {
        const res = await (0, _enquirer.prompt)([{
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
      const channel = await client.queryChannels({
        id: flags.channel,
        type: flags.type
      }, {
        last_message_at: -1
      }, {
        state: true
      });

      if (!channel.length) {
        this.log(`Channel ${_chalk.default.bold(flags.channel)} with type ${_chalk.default.bold(flags.type)} could not be found.`);
        this.exit();
      }

      this.log(JSON.stringify({ ...channel[0].data,
        members: channel[0].state.members
      }));
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

ChannelGet.flags = {
  channel: _command.flags.string({
    char: 'c',
    description: 'The channel ID you wish to retrieve.',
    required: false
  }),
  type: _command.flags.string({
    char: 't',
    description: 'Type of channel.',
    required: false
  })
};
ChannelGet.description = 'Gets a specific channel by its ID and type.';
module.exports.ChannelGet = ChannelGet;