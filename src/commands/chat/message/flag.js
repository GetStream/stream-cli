import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';

class MessageFlag extends Command {
	async run() {
		const { flags } = this.parse(MessageFlag);

		try {
			if (!flags.message) {
				const res = await prompt([
					{
						type: 'input',
						name: 'message',
						message: 'What is the unique identifier for the message?',
						required: true
					}
				]);

				flags.message = res.message;
			}

			const client = await chatAuth(this);
			const response = client.flagMessage(flags.message);

			if (flags.json) {
				this.log(JSON.stringify(response));
				this.exit();
			}

			this.log(`Message ${chalk.bold(flags.message)} has been flagged.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

MessageFlag.flags = {
	message: flags.string({
		char: 'm',
		description: 'The unique identifier of the message you want to flag.',
		required: false
	}),
	json: flags.boolean({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

MessageFlag.description = 'Flags a message.';

module.exports.MessageFlag = MessageFlag;
