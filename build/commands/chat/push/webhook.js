"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class PushWebhook extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(PushWebhook);

    try {
      if (!flags.url) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'url',
          message: `What is the absolute URL for your webhook?`,
          required: true
        }]);
        flags.url = res.url;
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      await client.updateAppSettings({
        webhook_url: flags.url
      });

      if (flags.json) {
        const settings = await client.getAppSettings();
        this.log(JSON.stringify({
          webhook_url: settings.app.webhook_url
        }));
        this.exit();
      }

      this.log(`Push notifications have been enabled for ${_chalk.default.bold('Webhooks')}.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

PushWebhook.flags = {
  url: _command.flags.string({
    char: 'u',
    description: 'A fully qualified URL for webhook support.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
PushWebhook.description = 'Tests webhook notifications.';
module.exports.PushWebhook = PushWebhook;