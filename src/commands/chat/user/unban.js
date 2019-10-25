const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class UserUnban extends Command {
	async run() {
		const { flags } = this.parse(UserUnban);

		try {
			if (!flags.user) {
				const res = await prompt([
					{
						type: 'input',
						name: 'user',
						message: 'What is the unique identifier for the user?',
						required: true,
					},
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await chatAuth(this);

			const unban = await client.unbanUser(flags.user);

			if (flags.json) {
				this.log(JSON.stringify(unban));
				this.exit();
			}

			this.log(`The user ${chalk.bold(flags.user)} has been unbanned.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

UserUnban.flags = {
	user: flags.string({
		char: 'u',
		description: 'The unique identifier of the user to unban.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.UserUnban = UserUnban;
