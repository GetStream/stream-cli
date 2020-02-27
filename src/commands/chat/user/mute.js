import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';

class UserMute extends Command {
	async run() {
		const { flags } = this.parse(UserMute);

		try {
			if (!flags.user) {
				const res = await prompt([
					{
						type: 'input',
						name: 'user',
						message: 'What is the unique identifier for the user?',
						required: true
					}
				]);

				flags.user = res.user;
			}

			const client = await chatAuth(this);
			const response = await client.muteUser(flags.user, 'CLI');

			if (flags.json) {
				this.log(JSON.stringify(response));
				this.exit();
			}

			this.log(`User ${chalk.bold(flags.user)} has been muted.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

UserMute.flags = {
	user: flags.string({
		char: 'u',
		description: 'The unique identifier of the user to mute.',
		required: false
	}),
	json: flags.boolean({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

UserMute.description = 'Mutes a user.';

module.exports.UserMute = UserMute;
