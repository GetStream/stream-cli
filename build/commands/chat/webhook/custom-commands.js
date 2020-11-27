"use strict";

var _command = require("@oclif/command");

var _webhook = require("../../../utils/webhook");

class CustomCommandsCommand extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(CustomCommandsCommand);
    await (0, _webhook.runWebhookCommand)("custom_action_handler_url", this, flags);
  }

}

(0, _webhook.setFlags)(CustomCommandsCommand, "Sets custom commands webhook URL");
module.exports.CustomCommandsCommand = CustomCommandsCommand;