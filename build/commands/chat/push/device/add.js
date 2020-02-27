"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chatAuth = require("../../../../utils/auth/chat-auth");

class DeviceAdd extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(DeviceAdd);

    try {
      if (!flags.user_id || !flags.provider || !flags.device_id) {
        const result = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'user_id',
          hint: 'user-123',
          message: 'What is the User ID?',
          required: true
        }, {
          type: 'input',
          name: 'device_id',
          hint: `device-123`,
          message: 'What is the Device ID?',
          required: true
        }, {
          type: 'select',
          name: 'provider',
          message: 'What is the push provider?',
          required: true,
          choices: [{
            message: 'APN',
            value: 'apn'
          }, {
            message: 'Firebase',
            value: 'firebase'
          }]
        }]);

        for (const key in result) {
          if (result.hasOwnProperty(key)) {
            flags[key] = result[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      await client.addDevice(flags.device_id || '', flags.provider || '', flags.user_id || '');
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

DeviceAdd.flags = {
  user_id: _command.flags.string({
    char: 'u',
    description: 'User ID',
    required: false
  }),
  device_id: _command.flags.string({
    char: 'd',
    description: 'Device id or token.',
    required: false
  }),
  provider: _command.flags.string({
    char: 'p',
    description: 'Push provider',
    required: false
  })
};
DeviceAdd.description = 'Adds a new device for push.';
module.exports.DeviceAdd = DeviceAdd;