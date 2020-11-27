"use strict";

var _command = require("@oclif/command");

var _webhook = require("../../../utils/webhook");

class PushCommand extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(PushCommand);
    await (0, _webhook.runWebhookCommand)("webhook_url", this, flags);
  }

}

(0, _webhook.setFlags)(PushCommand, "Sets push webhook URL");
module.exports.PushCommand = PushCommand;