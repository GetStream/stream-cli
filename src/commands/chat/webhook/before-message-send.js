import {Command} from "@oclif/command";

import { runWebhookCommand, setFlags} from 'utils/webhook';

class BeforeMessageSendCommand extends Command {
	async run() {
		const {flags} = this.parse(BeforeMessageSendCommand);
		await runWebhookCommand("before_message_send_hook_url", this, flags);
	}
}
setFlags(BeforeMessageSendCommand, "Sets before message send webhook URL")

module.exports.BeforeMessageSendCommand = BeforeMessageSendCommand
