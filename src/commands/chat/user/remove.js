const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class UserRemove extends Command {
	async run() {
		const { flags } = this.parse(UserRemove);

		try {
			if (
				!flags.channel ||
				!flags.type ||
				!flags.moderators ||
				!flags.json
			) {
				const res = await prompt([
					{
						type: 'input',
						name: 'channel',
						message: `What is the unique identifier for the channel?`,
						required: true,
					},
					{
						type: 'select',
						name: 'type',
						message: 'What type of channel is this?',
						required: true,
						choices: [
							{ message: 'Livestream', value: 'livestream' },
							{ message: 'Messaging', value: 'messaging' },
							{ message: 'Gaming', value: 'gaming' },
							{ message: 'Commerce', value: 'commerce' },
							{ message: 'Team', value: 'team' },
						],
					},
					{
						type: 'input',
						name: 'users',
						message: `What users would you like to remove (comma separated)?`,
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

			const channel = await client.channel(flags.type, flags.channel);
			const remove = await channel.demoteModerators(
				flags.users.split(',')
			);

			if (flags.json) {
				this.log(JSON.stringify(remove));
				this.exit();
			}

			this.log(
				`${chalk.bold(flags.users.length)} users have been removed.`
			);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

UserRemove.flags = {
	channel: flags.string({
		char: 'c',
		description: 'Channel name.',
		required: false,
	}),
	type: flags.string({
		char: 't',
		description: 'Channel type.',
		required: false,
	}),
	moderators: flags.string({
		char: 'm',
		description: 'Comma separated list of moderators to remove.',
		required: true,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.UserRemove = UserRemove;
