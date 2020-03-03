"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ChannelUpdateType extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(ChannelUpdateType);

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
          }, {
            message: 'Custom',
            value: 'custom'
          }]
        }, {
          type: 'select',
          name: 'typing_events',
          message: 'Would you like to enable typing events?',
          required: true,
          choices: [{
            message: 'Yes',
            value: true
          }, {
            message: 'No',
            value: false
          }]
        }, {
          type: 'select',
          name: 'read_events',
          message: 'Would you like to enable read events?',
          required: true,
          choices: [{
            message: 'Yes',
            value: true
          }, {
            message: 'No',
            value: false
          }]
        }, {
          type: 'select',
          name: 'read_events',
          message: 'Would you like to enable connect events?',
          required: true,
          choices: [{
            message: 'Yes',
            value: true
          }, {
            message: 'No',
            value: false
          }]
        }, {
          type: 'select',
          name: 'search',
          message: 'Would you like to enable search?',
          required: true,
          choices: [{
            message: 'Yes',
            value: true
          }, {
            message: 'No',
            value: false
          }]
        }, {
          type: 'select',
          name: 'reactions',
          message: 'Would you like to enable reactions?',
          required: true,
          choices: [{
            message: 'Yes',
            value: true
          }, {
            message: 'No',
            value: false
          }]
        }, {
          type: 'select',
          name: 'replies',
          message: 'Would you like to enable replies?',
          required: true,
          choices: [{
            message: 'Yes',
            value: true
          }, {
            message: 'No',
            value: false
          }]
        }, {
          type: 'select',
          name: 'mutes',
          message: 'Would you like to enable mutes?',
          required: true,
          choices: [{
            message: 'Yes',
            value: true
          }, {
            message: 'No',
            value: false
          }]
        }, {
          type: 'select',
          name: 'uploads',
          message: 'Would you like to enable uploads?',
          required: true,
          choices: [{
            message: 'Yes',
            value: true
          }, {
            message: 'No',
            value: false
          }]
        }, {
          type: 'select',
          name: 'url_enrichment',
          message: 'Would you like to enable uploads?',
          required: true,
          choices: [{
            message: 'Yes',
            value: true
          }, {
            message: 'No',
            value: false
          }]
        }, {
          type: 'select',
          name: 'automod',
          message: 'Would you like to enable automod?',
          required: true,
          choices: [{
            message: 'Simple',
            value: 'simple'
          }, {
            message: 'AI',
            value: 'AI'
          }, {
            message: 'Disabled',
            value: 'disabled'
          }]
        }, {
          type: 'select',
          name: 'message_retention',
          message: 'How many days would you like to retain your messages?',
          hint: 'infinite',
          required: false,
          choices: [{
            message: 'Infinite',
            value: 'infinite'
          }, {
            message: '30',
            value: '30'
          }, {
            message: '60',
            value: '60'
          }, {
            message: '90',
            value: '90'
          }, {
            message: '120',
            value: '120'
          }, {
            message: '150',
            value: '150'
          }, {
            message: '180',
            value: '180'
          }, {
            message: '210',
            value: '210'
          }, {
            message: '240',
            value: '240'
          }, {
            message: '270',
            value: '270'
          }, {
            message: '300',
            value: '300'
          }, {
            message: '330',
            value: '330'
          }, {
            message: '360',
            value: '360'
          }]
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }

        if (res.type === 'custom') {
          const channel = await (0, _enquirer.prompt)([{
            type: 'input',
            name: 'type',
            message: 'What is your custom channel type?',
            required: true
          }]);
          flags.type = channel.type;
        }
      }

      const payload = {};

      for (const key in flags) {
        if (flags.hasOwnProperty(key)) {
          if (flags[key] === 'true') flags[key] = true;
          if (flags[key] === 'false') flags[key] = false;
          if (flags[key] === 'No') flags[key] = false;
          if (/^\d+$/.test(flags[key])) flags[key] = parseInt(flags[key], 10);
          payload[key] = flags[key];
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const update = await client.updateChannelType(flags.type, { ...payload
      });

      if (flags.json) {
        this.log(JSON.stringify(update));
        this.exit();
      }

      this.log(`Channel type ${_chalk.default.bold(flags.type)} has been updated.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

ChannelUpdateType.flags = {
  type: _command.flags.string({
    char: 't',
    description: 'Type of channel.',
    required: false
  }),
  replies: _command.flags.boolean({
    char: 'r',
    description: 'Enable or disable replies (true/false)',
    required: false
  }),
  typing_events: _command.flags.boolean({
    char: 'y',
    description: 'Enable or disable typing events (true/false)',
    required: false
  }),
  read_events: _command.flags.string({
    char: 'e',
    description: 'Enable or disable read events (true/false)',
    required: false
  }),
  connect_events: _command.flags.boolean({
    char: 'c',
    description: 'Enable or disable connect events (true/false)',
    required: false
  }),
  search: _command.flags.boolean({
    char: 's',
    description: 'Enable or disable search (true/false)',
    required: false
  }),
  reactions: _command.flags.boolean({
    char: 'a',
    description: 'Enable or disable reactions (true/false)',
    required: false
  }),
  replies: _command.flags.boolean({
    char: 'p',
    description: 'Enable or disable replies (true/false)',
    required: false
  }),
  mutes: _command.flags.boolean({
    char: 'm',
    description: 'Enable or disable mutes (true/false)',
    required: false
  }),
  automod: _command.flags.string({
    char: 'a',
    description: 'Enable or disable automod (enabled/disabled)',
    required: false
  }),
  message_retention: _command.flags.string({
    char: 'a',
    description: 'How long to retain messages (defaults to infinite)',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
ChannelUpdateType.description = 'Updates a channels type configuration.';
module.exports.ChannelUpdateType = ChannelUpdateType;