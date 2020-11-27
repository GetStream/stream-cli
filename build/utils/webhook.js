"use strict";

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _chatAuth = require("./auth/chat-auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

module.exports.setFlags = (command, description) => {
  command.flags = {
    url: _command.flags.string({
      char: 'u',
      description: 'A fully qualified URL for webhook support.',
      required: false
    }),
    json: _command.flags.boolean({
      char: 'j',
      description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
      required: false
    })
  };
  command.description = description;
};

module.exports.runWebhookCommand = async (appFieldName, command, flags) => {
  try {
    if (!flags.url) {
      const res = await (0, _enquirer.prompt)([{
        type: 'input',
        name: 'url',
        message: `What is the absolute URL for your webhook?`,
        required: true
      }]);
      flags.url = res.url;
    }

    let settings = {
      [appFieldName]: flags.url
    };
    const client = await (0, _chatAuth.chatAuth)(command);
    await client.updateAppSettings(settings);

    if (flags.json) {
      const settings = await client.getAppSettings();
      command.log(JSON.stringify(settings));
      command.exit();
    }

    command.log(`Webhook ${_chalk.default.bold(appFieldName)} is set to ${flags.url}`);
    command.exit();
  } catch (error) {
    await command.config.runHook('telemetry', {
      ctx: command,
      error
    });
  }
};