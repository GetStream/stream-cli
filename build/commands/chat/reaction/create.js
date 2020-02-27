"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chatAuth = require("../../../utils/auth/chat-auth");

class ReactionCreate extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(ReactionCreate);

    try {
      if (!flags.json) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'channel',
          message: 'What is the unique identifier for the channel?',
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
          name: 'message',
          message: 'What is the unique identifier for the message?',
          required: true
        }, {
          type: 'input',
          name: 'reaction',
          hint: 'love',
          message: 'What is the reaction you would like to add?',
          required: true
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const channel = client.channel(flags.type, flags.channel);
      const reaction = await channel.sendReaction(flags.message, {
        type: flags.reaction
      });

      if (this.json) {
        this.log(JSON.stringify(reaction));
        this.exit();
      }

      this.log('Your reaction has been created.');
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

ReactionCreate.flags = {
  channel: _command.flags.string({
    char: 'c',
    description: 'The unique identifier for the channel.',
    required: false
  }),
  type: _command.flags.string({
    char: 't',
    description: 'The type of channel.',
    required: false
  }),
  message: _command.flags.string({
    char: 'c',
    description: 'The unique identifier for the message.',
    required: false
  }),
  reaction: _command.flags.string({
    char: 'r',
    description: 'A reaction for the message (e.g. love).',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
ReactionCreate.description = 'Creates a new reaction.';
module.exports.ReactionCreate = ReactionCreate;