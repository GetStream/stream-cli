import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';

class UserRemove extends Command {
	async run() {
		const { flags } = this.parse(UserRemove);

		try {
			if (!flags.user) {
				const res = await prompt([
					{
						type: 'input',
						name: 'user',
						message: 'What is the unique ID of the user you would like to remove?',
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

			const exists = await client.queryUsers({ id: flags.user }, { id: -1 });

			if (!exists.users.length) {
				this.log(`User ${flags.user} does not exist or has already been removed.`);
				this.exit();
			}

			const remove = await client.deleteUser(flags.user, {
				mark_messages_deleted: true,
				hard_delete: true
			});

			if (flags.json) {
				this.log(JSON.stringify(remove));
				this.exit();
			}

			this.log(`${chalk.bold(flags.user)} has been removed.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

UserRemove.flags = {
	user: flags.string({
		char: 'm',
		description: 'A unique ID of the user you would like to remove.',
		required: false
	}),
	json: flags.string({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

UserRemove.description = 'Allows for deactivating a user and wiping all of their messages.';

module.exports.UserRemove = UserRemove;
