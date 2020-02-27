"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _cliTable = _interopRequireDefault(require("cli-table"));

var _chatAuth = require("../../../../utils/auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class DeviceList extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(DeviceList);

    try {
      if (!flags.user_id) {
        const result = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'user_id',
          hint: 'user-123',
          message: 'What is the User ID?',
          required: true
        }]);

        for (const key in result) {
          if (result.hasOwnProperty(key)) {
            flags[key] = result[key];
          }
        }
      }

      const client = await (0, _chatAuth.chatAuth)(this);
      const response = await client.getDevices(flags.user_id || '');

      if (response.devices && response.devices.length !== 0) {
        const table = new _cliTable.default({
          head: ['Device ID', 'Push provider']
        });

        for (const device of response.devices) {
          table.push([device.id, device.push_provider]);
        }

        this.log(table.toString());
      } else {
        this.log('User has no devices');
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

DeviceList.flags = {
  user_id: _command.flags.string({
    char: 'u',
    description: 'User ID',
    required: false
  })
};
DeviceList.description = 'Gets all devices registered for push.';
module.exports.DeviceList = DeviceList;