"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _fsExtra = _interopRequireDefault(require("fs-extra"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ConfigDestroy extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(ConfigDestroy);

    const config = _path.default.join(this.config.configDir, 'config.json');

    try {
      if (!flags.force) {
        const answer = await (0, _enquirer.prompt)({
          type: 'confirm',
          name: 'continue',
          message: _chalk.default.red.bold(`This command will delete your current configuration. Are you sure you want to continue? ${_nodeEmoji.default.get('warning')} `)
        });

        if (answer.continue) {
          await _fsExtra.default.remove(config);
        }
      }

      this.log(`Config destroyed. Run the command ${_chalk.default.bold('stream config:set')} to generate a new stream configuration file.`);
      this.exit(0);
    } catch (error) {
      this.error(error, {
        exit: 1
      });
    }
  }

}

ConfigDestroy.flags = {
  force: _command.flags.boolean({
    char: 'f',
    description: 'Force remove Stream configuration from cache.',
    required: false
  })
};
ConfigDestroy.description = 'Destroys your user configuration.';
module.exports.ConfigDestroy = ConfigDestroy;