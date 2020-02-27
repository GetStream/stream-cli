"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class PushFirebase extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(PushFirebase);

    try {
      const client = await (0, _chatAuth.chatAuth)(this);

      if (flags.disable) {
        const result = await (0, _enquirer.prompt)({
          type: 'toggle',
          name: 'proceed',
          message: 'This will disable Firebase push notifications and remove your Firebase Server Key. Are you sure?',
          required: true
        });

        if (result.proceed) {
          await client.updateAppSettings({
            firebase_config: {
              disabled: true
            }
          });
          this.log(`Push notifications have been ${_chalk.default.red('disabled')} with ${_chalk.default.bold('Firebase')}.`);
        }

        this.exit();
      } else if (!flags.key) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'key',
          message: `What is your Server key for Firebase?`,
          required: true
        }, {
          type: 'input',
          name: 'notification_template',
          hint: 'Omit for Stream default',
          message: `What JSON notification template would you like to use?`,
          required: false
        }, {
          type: 'input',
          name: 'data_template',
          hint: 'Omit for Stream default',
          message: `What JSON data template would you like to use?`,
          required: false
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      const payload = {
        firebase_config: {
          server_key: flags.key
        }
      };

      if (flags.notification_template) {
        payload.firebase_config.notification_template = flags.notification_template;
      }

      if (flags.data_template) {
        payload.firebase_config.data_template = flags.data_template;
      }

      await client.updateAppSettings(payload);

      if (flags.json) {
        const settings = await client.getAppSettings();
        this.log(JSON.stringify(settings.app.push_notifications.firebase));
        this.exit();
      }

      this.log(`Push notifications have been ${_chalk.default.green('enabled')} for ${_chalk.default.bold('Firebase')}.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

PushFirebase.flags = {
  key: _command.flags.string({
    char: 'k',
    description: 'Server key for Firebase.',
    required: false
  }),
  notification_template: _command.flags.string({
    char: 'n',
    description: 'JSON notification template.',
    required: false
  }),
  data_template: _command.flags.string({
    char: 'd',
    description: 'JSON data template.',
    required: false
  }),
  disable: _command.flags.boolean({
    description: 'Disable Firebase push notifications and clear config.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
PushFirebase.description = 'Specifies Firebase for push notifications.';
module.exports.PushFirebase = PushFirebase;