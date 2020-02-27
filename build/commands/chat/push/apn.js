"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _fs = _interopRequireDefault(require("fs"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class PushApn extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(PushApn);

    try {
      const client = await (0, _chatAuth.chatAuth)(this);

      if (flags.disable) {
        const result = await (0, _enquirer.prompt)({
          type: 'toggle',
          name: 'proceed',
          message: 'This will disable APN push notifications and remove your APN settings, such as certificates or tokens. Are you sure?',
          required: true
        });

        if (result.proceed) {
          await client.updateAppSettings({
            apn_config: {
              disabled: true
            }
          });
          this.log(`Push notifications have been ${_chalk.default.red('disabled')} with ${_chalk.default.bold('APN')}.`);
        }

        this.exit();
      } else if (!flags.p12_cert && !flags.auth_key) {
        const type = await (0, _enquirer.prompt)([{
          type: 'select',
          name: 'apn',
          message: 'What type of APN configuration will you use?',
          required: true,
          choices: [{
            message: 'certificate',
            value: 'cert'
          }, {
            message: 'token',
            value: 'token'
          }]
        }]);
        let res = {};

        if (type.apn === 'cert') {
          res = await (0, _enquirer.prompt)([{
            type: 'input',
            name: 'p12_cert',
            hint: '/Users/username/Desktop/cert.p12',
            message: 'What is the absolute path to your .p12 file?',
            required: true
          }, {
            type: 'input',
            name: 'notification_template',
            hint: 'Omit for Stream default',
            message: `What JSON notification template would you like to use?`,
            required: false
          }]);
        }

        if (type.apn === 'token') {
          res = await (0, _enquirer.prompt)([{
            type: 'input',
            name: 'auth_key',
            hint: '/Users/username/Desktop/key.p8',
            message: 'What is the filepath to your .p8 file?',
            required: true
          }, {
            type: 'input',
            name: 'key_id',
            message: 'What is your key ID?',
            required: true
          }, {
            type: 'input',
            name: 'team_id',
            message: 'What is your team ID?',
            required: true
          }, {
            type: 'input',
            name: 'notification_template',
            hint: 'Omit for Stream default',
            message: `What JSON notification template would you like to use?`,
            required: false
          }]);
        }

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }

        res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'bundle_id',
          hint: 'com.apple.test',
          message: 'What is your bundle ID?',
          required: true
        }, {
          type: 'toggle',
          name: 'development',
          message: 'Do you want to use the development endpoint for APNs?',
          required: true
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      if (flags.p12_cert && flags.auth_key) {
        this.error(`Can't have both certificate and auth key`, {
          exit: 1
        });
      }

      const payload = {
        apn_config: {
          p12_cert: _fs.default.existsSync(flags.p12_cert) ? _fs.default.readFileSync(flags.p12_cert) : '',
          auth_key: _fs.default.existsSync(flags.auth_key) ? _fs.default.readFileSync(flags.auth_key, 'utf-8') : '',
          key_id: flags.key_id || '',
          team_id: flags.team_id || '',
          bundle_id: flags.bundle_id || '',
          development: flags.development || false
        }
      };

      if (flags.notification_template) {
        payload.apn_config.notification_template = flags.notification_template;
      }

      let auth_type = '';

      if (flags.p12_cert) {
        auth_type = 'certificate';
      } else {
        auth_type = 'token';
      }

      payload.apn_config.auth_type = auth_type;
      await client.updateAppSettings(payload);

      if (flags.json) {
        const settings = await client.getAppSettings();
        this.log(JSON.stringify(settings.app.push_notifications));
        this.exit();
      }

      this.log(`Push notifications have been ${_chalk.default.green('enabled')} with ${_chalk.default.bold('APN')}.`);
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

PushApn.flags = {
  auth_key: _command.flags.string({
    char: 'a',
    description: 'Absolute path to .p8 auth key.',
    required: false
  }),
  key_id: _command.flags.string({
    char: 'k',
    description: 'Key ID.',
    required: false
  }),
  bundle_id: _command.flags.string({
    char: 'b',
    description: 'Bundle identifier (e.g. com.apple.test).',
    required: false
  }),
  team_id: _command.flags.string({
    char: 't',
    description: 'Team ID.',
    required: false
  }),
  development: _command.flags.boolean({
    char: 'd',
    description: 'Use development endpoint',
    required: false
  }),
  p12_cert: _command.flags.string({
    char: 'c',
    description: 'Absolute path to .p12 file.',
    required: false
  }),
  notification_template: _command.flags.string({
    char: 'n',
    description: 'JSON template for notifications.',
    required: false
  }),
  disable: _command.flags.boolean({
    description: 'Disable APN push notifications and clear config.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
PushApn.description = 'Specifies APN for push notifications.';
module.exports.PushApn = PushApn;