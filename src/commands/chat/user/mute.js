const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { auth } = require('../../../utils/auth');

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
						required: true,
					},
				]);

				flags.user = res.user;
			}

			const client = await auth(this);
			const flag = client.muteUser(flags.user);

			if (flags.json) {
				this.log(JSON.stringify(flag));
				this.exit();
			}

			this.log(`User ${chalk.bold(flags.user)} has been muted.`);
			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

UserMute.flags = {
	user: flags.string({
		char: 'u',
		description: 'The unique identifier of the user to mute.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.UserMute = UserMute;
