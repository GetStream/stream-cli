"use strict";

var _command = require("@oclif/command");

var _webhook = require("../../../utils/webhook");

class BeforeMessageSendCommand extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(BeforeMessageSendCommand);
    await (0, _webhook.runWebhookCommand)("before_message_send_hook_url", this, flags);
  }

}

(0, _webhook.setFlags)(BeforeMessageSendCommand, "Sets before message send webhook URL");
module.exports.BeforeMessageSendCommand = BeforeMessageSendCommand;