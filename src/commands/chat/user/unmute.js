import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';

class UserUnmute extends Command {
	async run() {
		const { flags } = this.parse(UserUnmute);

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
			const response = await client.unmuteUser(flags.user);

			if (flags.json) {
				this.log(JSON.stringify(response));
				this.exit();
			}

			this.log(`User ${chalk.bold(flags.user)} has been unmuted.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

UserUnmute.flags = {
	user: flags.string({
		char: 'u',
		description: 'The unique identifier of the user to unmute.',
		required: false
	}),
	json: flags.boolean({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

UserUnmute.description = 'Unmutes a user.';

module.exports.UserUnmute = UserUnmute;
