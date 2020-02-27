"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chatAuth = require("../../../../utils/auth/chat-auth");

class DeviceDelete extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(DeviceDelete);

    try {
      if (!flags.user_id || !flags.device_id) {
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
        }]);

        for (const key in result) {
          if (result.hasOwnProperty(key)) {
            flags[key] = result[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      await client.removeDevice(flags.device_id || '', flags.user_id || '');
      this.exit();
    } catch (error) {
      await this.config.runHook('telemetry', {
        ctx: this,
        error
      });
    }
  }

}

DeviceDelete.flags = {
  user_id: _command.flags.string({
    char: 'u',
    description: 'User ID',
    required: false
  }),
  device_id: _command.flags.string({
    char: 'd',
    description: 'Device id or token.',
    required: false
  })
};
DeviceDelete.description = 'Removes a device from push.';
module.exports.DeviceDelete = DeviceDelete;