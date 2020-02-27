import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';

class UserDeactivate extends Command {
	async run() {
		const { flags } = this.parse(UserDeactivate);

		try {
			if (!flags.user || !flags.hard) {
				const res = await prompt([
					{
						type: 'input',
						name: 'user',
						message: 'What is the unique ID of the user you would like to deactivate?',
						required: true
					},
					{
						type: 'select',
						name: 'hard',
						message: 'Would you like to perform a hard delete on messages?',
						required: true,
						choices: [
							{
								message: 'No',
								value: false
							},
							{
								message: 'Yes',
								value: true
							}
						]
					}
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await chatAuth(this);

			const exists = await client.queryUsers({ id: flags.user }, { id: -1 });

			if (!exists.users.length) {
				this.log(`User ${flags.user} does not exist or has already been deactivated.`);
				this.exit();
			}

			const deactivate = await client.deactivateUser(flags.user, {
				mark_messages_deleted: Boolean(flags.hard)
			});

			if (flags.json) {
				this.log(JSON.stringify(deactivate));
				this.exit();
			}

			this.log(`${chalk.bold(flags.user)} has been deactivated.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

UserDeactivate.flags = {
	user: flags.string({
		char: 'm',
		description: 'A unique ID of the user you would like to deactivate.',
		required: false
	}),
	hard: flags.string({
		char: 'h',
		description: 'Hard deletes all messages associated with the user.',
		required: false
	}),
	json: flags.string({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

UserDeactivate.description = 'Allows for deactivating a user and wiping all of their messages.';

module.exports.UserDeactivate = UserDeactivate;
