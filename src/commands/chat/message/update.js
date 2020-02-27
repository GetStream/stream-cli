import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';
import { credentials } from 'utils/config';

class MessageUpdate extends Command {
	async run() {
		const { flags } = this.parse(MessageUpdate);

		try {
			const { name } = await credentials(this);

			if (!flags.message || !flags.text) {
				const res = await prompt([
					{
						type: 'input',
						name: 'message',
						message: `What is the unique identifier for the message?`,
						required: true
					},
					{
						type: 'input',
						name: 'text',
						message: 'What is the updated message?',
						required: true
					}
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const payload = {
				id: flags.message,
				text: flags.text,
				user: {
					id: 'CLI',
					name
				}
			};

			if (flags.attachments) {
				payload.attachments = JSON.parse(flags.attachments);
			}

			const client = await chatAuth(this);

			await client.setUser({
				id: 'CLI',
				status: 'invisible'
			});

			const update = await client.updateMessage(payload);

			if (flags.json) {
				this.log(JSON.stringify(update));
				this.exit();
			}

			this.log(`Message ${chalk.bold(flags.message)} has been updated.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

MessageUpdate.flags = {
	message: flags.string({
		char: 'm',
		description: 'The unique identifier for the message.',
		required: false
	}),
	text: flags.string({
		char: 't',
		description: 'The message you would like to send as text.',
		required: false
	}),
	attachments: flags.string({
		char: 'a',
		description: 'A JSON payload of attachments to send along with a message.',
		required: false
	}),
	json: flags.boolean({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

MessageUpdate.description = 'Updates a message.';

module.exports.MessageUpdate = MessageUpdate;
