const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class UserReactivate extends Command {
	async run() {
		const { flags } = this.parse(UserReactivate);

		this.log(flags);

		try {
			if (!flags.user || !flags.restore) {
				const res = await prompt([
					{
						type: 'input',
						name: 'user',
						message:
							'What is the unique ID of the user you would like to reactivate?',
						required: true,
					},
					{
						type: 'select',
						name: 'restore',
						message: 'Would you like to restore all messages?',
						required: true,
						choices: [
							{
								message: 'No',
								value: false,
							},
							{
								message: 'Yes',
								value: true,
							},
						],
					},
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await chatAuth(this);

			const { user } = await client.reactivateUser(flags.user, {
				restore_messages: Boolean(flags.restore),
			});

			if (flags.json) {
				this.log(JSON.stringify(user));
				this.exit();
			}

			this.log(`${chalk.bold(flags.user)} has been reactivated.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

UserReactivate.flags = {
	user: flags.string({
		char: 'm',
		description: 'A unique ID of the user you would like to reactivate.',
		required: false,
	}),
	restore: flags.string({
		char: 'r',
		description: 'Restores all deleted messages associated with the user.',
		required: false,
	}),
	json: flags.string({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

UserReactivate.description =
	'Reactivates a user who was previously deactivated.';

module.exports.UserReactivate = UserReactivate;
