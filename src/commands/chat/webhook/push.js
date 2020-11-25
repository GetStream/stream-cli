import {Command} from "@oclif/command";

import { runWebhookCommand, setFlags} from 'utils/webhook';

class PushCommand extends Command {
	async run() {
		const {flags} = this.parse(PushCommand);
		await runWebhookCommand("webhook_url", this, flags);
	}
}
setFlags(PushCommand, "Sets push webhook URL")

module.exports.PushCommand = PushCommand
