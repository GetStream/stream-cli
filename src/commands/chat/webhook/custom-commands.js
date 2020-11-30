import {Command} from "@oclif/command";

import { runWebhookCommand, setFlags} from 'utils/webhook';

class CustomCommandsCommand extends Command {
	async run() {
		const {flags} = this.parse(CustomCommandsCommand);
		await runWebhookCommand("custom_action_handler_url", this, flags);
	}
}
setFlags(CustomCommandsCommand, "Sets custom commands webhook URL")

module.exports.CustomCommandsCommand = CustomCommandsCommand
