"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _cliTable = _interopRequireDefault(require("cli-table"));

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class PushTest extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(PushTest);

    try {
      if (!flags.user_id) {
        const result = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'user_id',
          hint: 'user-123',
          message: 'What is the ID of the user you wish to test push settings with?',
          required: true
        }, {
          type: 'input',
          name: 'message_id',
          hint: `Omit to pick random message`,
          message: 'What is the message ID you want to use to render the notification template with?',
          required: false
        }, {
          type: 'input',
          name: 'apn_notification_template',
          hint: 'Omit for the APN template configured in your app',
          message: `What JSON notification template would you like to use for APN?`,
          required: false
        }, {
          type: 'input',
          name: 'firebase_notification_template',
          hint: 'Omit for the Firebase template configured in your app',
          message: `What JSON notification template would you like to use for Firebase?`,
          required: false
        }, {
          type: 'input',
          name: 'firebase_data_template',
          hint: 'Omit for the Firebase data template configured in your app',
          message: `What JSON data template would you like to use for Firebase?`,
          required: false
        }, {
          type: 'input',
          name: 'skip_devices',
          hint: 'Set if you want to skip sending to devices',
          message: `Do you want to skip sending to devices?`,
          default: false,
          required: false
        }]);

        for (const key in result) {
          if (result.hasOwnProperty(key)) {
            flags[key] = result[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const payload = {
        messageID: flags.message_id || '',
        apnTemplate: flags.apn_notification_template || '',
        firebaseTemplate: flags.firebase_notification_template || '',
        firebaseDataTemplate: flags.firebase_data_template || '',
        skipDevices: flags.skip_devices || false
      };
      const userID = flags.user_id || '';
      const response = await client.testPushSettings(userID, payload);

      if (flags.json) {
        this.log(JSON.stringify(response));
        this.exit();
      }

      if (response.general_errors) {
        this.log(`It seems there were some ${_chalk.default.red.bold('errors')} with the input you provided. Listing them below:`);

        for (const err of response.general_errors) {
          this.log(`\t - ${err}`);
        }
      }

      if (response.rendered_apn_template) {
        this.log(`Here is the rendered APN notification that will be sent to your devices:`);
        this.log(JSON.stringify(JSON.parse(response.rendered_apn_template), null, 4));
      }

      if (response.rendered_firebase_template) {
        this.log(`Here is the rendered Firebase notification that will be sent to your devices:`);
        this.log(JSON.stringify(JSON.parse(response.rendered_firebase_template), null, 4));
      }

      if (response.rendered_message) {
        this.log(`Here is the rendered notification payload that will be sent to your devices:`);
        this.log(JSON.stringify(JSON.parse(response.rendered_message), null, 4));
      }

      if (response.device_errors) {
        this.log(`It seems we couldn't push the notification to all of the user devices. Here are the ${_chalk.default.red.bold('errors')} for each device:`);
        const table = new _cliTable.default({
          head: ['Device ID', 'Push provider', 'Error message']
        });

        for (const [deviceID, errorDetails] of Object.entries(response.device_errors)) {
          table.push([deviceID || 'N/A', errorDetails.provider || 'N/A', errorDetails.error_message || 'N/A']);
        }

        this.log(table.toString());
      }

      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

PushTest.flags = {
  user_id: _command.flags.string({
    char: 'u',
    description: 'User ID',
    required: false
  }),
  message_id: _command.flags.string({
    char: 'm',
    description: 'Message ID.',
    required: false
  }),
  apn_notification_template: _command.flags.string({
    char: 'a',
    description: 'APN notification template',
    required: false
  }),
  firebase_notification_template: _command.flags.string({
    char: 'f',
    description: 'Firebase notification template',
    required: false
  }),
  firebase_data_template: _command.flags.string({
    char: 'd',
    description: 'Firebase data template',
    required: false
  }),
  skip_devices: _command.flags.string({
    char: 's',
    description: 'Skip devices',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
PushTest.description = 'Tests push notifications.';
module.exports.PushTest = PushTest;