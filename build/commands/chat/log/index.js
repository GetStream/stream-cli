"use strict";

var _command = require("@oclif/command");

var _jsonStringifyPrettyCompact = _interopRequireDefault(require("json-stringify-pretty-compact"));

var _enquirer = require("enquirer");

var _cardinal = _interopRequireDefault(require("cardinal"));

var _moment = _interopRequireDefault(require("moment"));

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class Log extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(Log);

    try {
      if (!flags.channel || !flags.type || !flags.event) {
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
          type: 'select',
          name: 'event',
          message: 'What event would you like to filter on?',
          required: true,
          choices: [{
            message: 'All Events - JSON',
            value: 'all'
          }, {
            message: 'User Status - Changed',
            value: 'user.status.changed'
          }, {
            message: 'User Watching - Start',
            value: 'user.watching.start'
          }, {
            message: 'User Watching - Stop',
            value: 'user.watching.stop'
          }, {
            message: 'User Updated',
            value: 'user.updated'
          }, {
            message: 'Typing - Start',
            value: 'typing.start'
          }, {
            message: 'Typing - Stop',
            value: 'typing.stop'
          }, {
            message: 'Message - New',
            value: 'message.new'
          }, {
            message: 'Message - Updated',
            value: 'message.updated'
          }, {
            message: 'Message - Deleted',
            value: 'message.deleted'
          }, {
            message: 'Message - Seen',
            value: 'message.seen'
          }, {
            message: 'Message - Reaction',
            value: 'message.reaction'
          }, {
            message: 'Member - Added',
            value: 'member.added'
          }, {
            message: 'Member - Removed',
            value: 'member.removed'
          }, {
            message: 'Channel - Updated',
            value: 'channel.updated'
          }, {
            message: 'Health - Check',
            value: 'health.check'
          }, {
            message: 'Connection - Changed',
            value: 'connection.changed'
          }, {
            message: 'Connection - Recovered',
            value: 'connection.recovered'
          }]
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      await client.setUser({
        id: 'CLI',
        role: 'admin',
        status: 'invisible'
      });
      const channel = client.channel(flags.type, flags.channel);
      await channel.watch();
      const format = 'dddd, MMMM Do YYYY [at] h:mm:ss A';

      if (flags.json) {
        channel.on(event => {
          this.log(JSON.stringify(event));
        });
      } else if (flags.event === 'all') {
        this.log(`Logging real-time events for ${flags.event}...`);
        channel.on(event => {
          const timestamp = _chalk.default.bold.green((0, _moment.default)(event.created_at).format(format));

          const payload = `${timestamp}: ${_chalk.default.bold(event.user?.name || event.user?.id)} performed event ${_chalk.default.bold(event.type)} in channel ${_chalk.default.bold(flags.channel)}.`;
          this.log(payload);
        });
      } else {
        this.log(`Logging real-time events for ${flags.event}...`);
        channel.on(flags.event, event => {
          const timestamp = _chalk.default.bold((0, _moment.default)(event.created_at).format(format));

          const payload = _cardinal.default.highlight((0, _jsonStringifyPrettyCompact.default)(event, {
            maxLength: 100
          }), {
            linenos: true
          });

          this.log(`${timestamp}:`, '\n\n', payload, '\n\n');
        });
      }
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

Log.flags = {
  channel: _command.flags.string({
    char: 'c',
    description: 'The channel ID you wish to log.',
    required: false
  }),
  type: _command.flags.string({
    char: 't',
    description: 'The type of channel.',
    required: false
  }),
  event: _command.flags.string({
    char: 'e',
    description: 'The type of event you want to listen on.',
    options: ['all', 'user.status.changed', 'user.watching.start', 'user.watching.stop', 'user.updated', 'typing.start', 'typing.stop', 'message.new', 'message.updated', 'message.deleted', 'message.seen', 'message.reaction', 'member.added', 'member.removed', 'channel.updated', 'health.check', 'connection.changed', 'connection.recovered'],
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
Log.description = 'Logs events in realtime.';
module.exports.Log = Log;