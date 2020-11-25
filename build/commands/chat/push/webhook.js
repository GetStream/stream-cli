"use strict";

var _command = require("@oclif/command");

var _webhook = require("../../../utils/webhook");

class WebhookCommand extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(WebhookCommand);
    await (0, _webhook.runWebhookCommand)("webhook_url", this, flags);
  }

}

(0, _webhook.setFlags)(WebhookCommand, "Sets push webhook URL");
module.exports.WebhookCommand = WebhookCommand;