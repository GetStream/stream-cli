"use strict";

var _command = require("@oclif/command");

var _cliTable = _interopRequireDefault(require("cli-table"));

var _chalk = _interopRequireDefault(require("chalk"));

var _config = require("../../utils/config");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ConfigGet extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(ConfigGet);

    try {
      const {
        name,
        email,
        apiKey,
        apiSecret,
        environment,
        telemetry,
        timeout
      } = await (0, _config.credentials)(this);

      if (flags.json) {
        this.log(JSON.stringify(await (0, _config.credentials)(this)));
        this.exit(0);
      }

      const table = new _cliTable.default();
      table.push({
        [`${_chalk.default.green.bold('Name')}`]: name
      }, {
        [`${_chalk.default.green.bold('Email')}`]: email
      }, {
        [`${_chalk.default.green.bold('API Key')}`]: apiKey
      }, {
        [`${_chalk.default.green.bold('API Secret')}`]: apiSecret
      }, {
        [`${_chalk.default.green.bold('Environment')}`]: environment
      }, {
        [`${_chalk.default.green.bold('Telemetry')}`]: telemetry
      }, {
        [`${_chalk.default.green.bold('Timeout(ms)')}`]: timeout
      });
      this.log(table.toString());
      this.exit(0);
    } catch (error) {
      this.error(error || 'A Stream CLI error has occurred.', {
        exit: 1
      });
    }
  }

}

ConfigGet.flags = {
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
ConfigGet.description = 'Outputs your user configuration.';
module.exports.ConfigGet = ConfigGet;