"use strict";

var _command = require("@oclif/command");

var _cliTable = _interopRequireDefault(require("cli-table"));

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class PushGet extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(PushGet);

    try {
      const client = await (0, _chatAuth.chatAuth)(this);
      const settings = await client.getAppSettings();

      if (flags.json) {
        this.log(JSON.stringify(settings.app));
        this.exit();
      }

      const table = new _cliTable.default();
      table.push({
        [`${_chalk.default.green.bold('APN')}`]: settings.app.push_notifications.apn.enabled ? 'Enabled' : 'Disabled'
      }, {
        [`${_chalk.default.green.bold('APN - Host')}`]: !settings.app.push_notifications.apn.host ? 'N/A' : settings.app.push_notifications.apn.host
      }, {
        [`${_chalk.default.green.bold('APN – Auth Type')}`]: !settings.app.push_notifications.apn.auth_type ? 'N/A' : settings.app.push_notifications.apn.auth_type.toUpperCase()
      }, {
        [`${_chalk.default.green.bold('APN – Key ID')}`]: !settings.app.push_notifications.apn.key_id ? 'N/A' : settings.app.push_notifications.apn.key_id
      }, {
        [`${_chalk.default.green.bold('APN – Notification Template')}`]: !settings.app.push_notifications.apn.notification_template ? 'Stream Default' : settings.app.push_notifications.apn.notification_template
      }, {
        [`${_chalk.default.green.bold('Firebase')}`]: settings.app.push_notifications.firebase.enabled ? 'Enabled' : 'Disabled'
      }, {
        [`${_chalk.default.green.bold('Firebase – Notification Template')}`]: !settings.app.push_notifications.firebase.notification_template ? 'Stream Default' : settings.app.push_notifications.firebase.notification_template
      }, {
        [`${_chalk.default.green.bold('Firebase – Data Template')}`]: !settings.app.push_notifications.firebase.data_template ? 'Stream Default' : settings.app.push_notifications.firebase.data_template
      }, {
        [`${_chalk.default.green.bold('Webhook – URL')}`]: !settings.app.webhook_url ? 'N/A' : settings.app.webhook_url
      });
      this.log(table.toString());
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

PushGet.flags = {
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
PushGet.description = 'Gets push notification settings.';
module.exports.PushGet = PushGet;