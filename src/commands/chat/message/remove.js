import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';

class MessageRemove extends Command {
	async run() {
		const { flags } = this.parse(MessageRemove);

		try {
			if (!flags.message) {
				const res = await prompt([
					{
						type: 'input',
						name: 'message',
						message: `What is the unique identifier for the message?`,
						required: true
					}
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await chatAuth(this);
			const remove = await client.deleteMessage(flags.message);

			if (flags.json) {
				this.log(JSON.stringify(remove));
				this.exit();
			}

			this.log(`The message ${chalk.bold(flags.message)} has been removed.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

MessageRemove.flags = {
	message: flags.string({
		char: 'message',
		description: 'The unique identifier of the message you would like to remove.',
		required: false
	}),
	json: flags.boolean({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

MessageRemove.description = 'Removes a message.';

module.exports.MessageRemove = MessageRemove;
