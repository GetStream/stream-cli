"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

var _config = require("../../../utils/config");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ChannelUpdate extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(ChannelUpdate);

    try {
      const {
        name
      } = await (0, _config.credentials)(this);

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
        }, {
          type: 'input',
          name: 'name',
          message: `What is the new name for the channel?`,
          required: false
        }, {
          type: 'input',
          name: 'image',
          message: `What is the absolute image URL for the channel?`,
          required: false
        }, {
          type: 'input',
          name: 'description',
          message: `What description would you like to set for the channel?`,
          required: false
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const channel = await client.channel(flags.type, flags.channel);
      const payload = {
        updated_by: {
          id: 'CLI',
          name
        }
      };
      if (flags.name) payload.name = flags.name;
      if (flags.image) payload.image = flags.image;
      if (flags.description) payload.description = flags.description;
      const update = await channel.update(payload);

      if (flags.json) {
        this.log(JSON.stringify(update));
        this.exit();
      }

      this.log(`Channel ${_chalk.default.bold(flags.channel)} has been updated.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

ChannelUpdate.flags = {
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
ChannelUpdate.description = 'Updates a channel.';
module.exports.ChannelUpdate = ChannelUpdate;