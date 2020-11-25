import {Command} from "@oclif/command";

import { runWebhookCommand, setFlags} from 'utils/webhook';

class WebhookCommand extends Command {
	async run() {
		const {flags} = this.parse(WebhookCommand);
		await runWebhookCommand("webhook_url", this, flags);
	}
}
setFlags(WebhookCommand, "Sets push webhook URL")

module.exports.WebhookCommand = WebhookCommand
