"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _path = _interopRequireDefault(require("path"));

var _fsExtra = _interopRequireDefault(require("fs-extra"));

var _config = require("../../utils/config");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ConfigSet extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(ConfigSet);

    const config = _path.default.join(this.config.configDir, 'config.json');

    try {
      if (!flags.name || !flags.email || !flags.key || !flags.secret || !flags.url || !flags.environment) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'name',
          message: `What is your full name?`,
          required: true
        }, {
          type: 'input',
          name: 'email',
          message: `What is your email address associated with Stream?`,
          required: true
        }, {
          type: 'input',
          name: 'key',
          message: `What is your Stream API key?`,
          required: true
        }, {
          type: 'password',
          name: 'secret',
          message: `What is your Stream API secret?`,
          required: true
        }, {
          type: 'input',
          name: 'environment',
          message: `What environment would you like to run in?`,
          default: 'production',
          required: false
        }, {
          type: 'input',
          name: 'telemetry',
          message: `Would you like to enable error tracking for debugging purposes?`,
          default: true,
          required: false
        }, {
          type: 'input',
          name: 'timeout',
          message: `Do you want to set a different timeout for requests?`,
          default: 3000,
          required: false
        }]);

        for (const key in res) {
          if (res.hasOwnProperty(key)) {
            flags[key] = res[key];
          }
        }
      }

      await _fsExtra.default.ensureDir(this.config.configDir);
      await _fsExtra.default.writeJson(config, {
        name: flags.name,
        email: flags.email.toLowerCase(),
        apiKey: flags.key,
        apiSecret: flags.secret,
        environment: flags.environment,
        telemetry: flags.telemetry !== undefined ? flags.telemetry : true,
        timeout: flags.timeout ? flags.timeout : 3000
      });

      if (flags.json) {
        this.log(JSON.stringify(await (0, _config.credentials)(this)));
        this.exit();
      }

      this.log('Your Stream CLI configuration has been generated!', _nodeEmoji.default.get('rocket'));
      this.exit();
    } catch (error) {
      this.error(error || 'A Stream CLI error has occurred.', {
        exit: 1
      });
    }
  }

}

ConfigSet.flags = {
  name: _command.flags.string({
    char: 'n',
    description: 'Full name for configuration.',
    required: false
  }),
  email: _command.flags.string({
    char: 'e',
    description: 'Email for configuration.',
    required: false
  }),
  key: _command.flags.string({
    char: 'k',
    description: 'API key for configuration.',
    required: false
  }),
  secret: _command.flags.string({
    char: 's',
    description: 'API secret for configuration.',
    required: false
  }),
  url: _command.flags.string({
    char: 'u',
    description: 'API base URL for configuration.',
    required: false
  }),
  environment: _command.flags.string({
    char: 'v',
    description: 'Environment to run in (production or development for token and permission checking).',
    required: false
  }),
  telemetry: _command.flags.boolean({
    char: 't',
    description: 'Enable error reporting for debugging purposes.',
    required: false
  }),
  timeout: _command.flags.integer({
    char: 'o',
    description: 'Timeout for requests in ms.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
ConfigSet.description = 'Sets your user configuration.';
module.exports.ConfigSet = ConfigSet;