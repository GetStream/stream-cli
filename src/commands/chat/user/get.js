const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class UserGet extends Command {
	async run() {
		const { flags } = this.parse(UserGet);

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
			const user = await client.queryUsers(
				{ id: { $in: [flags.user] } },
				{ id: -1 }
			);

			if (!user.users.length) {
				this.log(`User ${chalk.bold(flags.user)} could not be found.`);
				this.exit();
			}

			if (flags.json) {
				this.log(JSON.stringify(user.users[0]));
				this.exit();
			}

			this.log(user.users[0]);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

UserGet.flags = {
	user: flags.string({
		char: 'u',
		description: 'The unique identifier of the user to get.',
		required: false,
	}),
	presence: flags.string({
		char: 'p',
		description: 'Display the current status of the user.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.UserGet = UserGet;
