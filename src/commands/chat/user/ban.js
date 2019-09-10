const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class UserBan extends Command {
	async run() {
		const { flags } = this.parse(UserBan);

		try {
			let type;
			if (!flags.type) {
				type = await prompt([
					{
						type: 'select',
						name: 'ban',
						message:
							'Would you like to apply a global or channel ban?',
						required: true,
						choices: [
							{ message: 'Global', value: 'global' },
							{ message: 'Channel', value: 'channel' },
						],
					},
				]);
			}

			let cid;
			if (type.ban === 'channel') {
				cid = await prompt([
					{
						type: 'input',
						name: 'id',
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
				]);
			}

			if (!flags.user || !flags.reason || !flags.duration) {
				const res = await prompt([
					{
						type: 'input',
						name: 'user',
						message: 'What is the unique identifier for the user?',
						required: true,
					},
					{
						type: 'input',
						name: 'reason',
						message: 'What is the reason for banning the user?',
						required: true,
					},
					{
						type: 'input',
						name: 'duration',
						hint: 'minutes',
						message: 'How long would you like to ban the user for?',
						required: false,
					},
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await chatAuth(this);
			const payload = {
				reason: flags.reason,
				user_id: 'CLI',
			};

			if (flags.duration) {
				payload.timeout = parseInt(flags.duration * 60, 10);
			}

			let ban;
			if (type.ban === 'channel') {
				ban = await client
					.channel(cid.type, cid.id)
					.banUser(flags.user, payload);
			}

			if (type.ban === 'global') {
				ban = await client.banUser(flags.user, payload);
			}

			if (flags.json) {
				this.log(JSON.stringify(ban));
				this.exit();
			}

			this.log(`The user ${chalk.bold(flags.user)} has been banned.`);
			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

UserBan.flags = {
	type: flags.string({
		char: 'u',
		description: 'Type of ban to perform (e.g. global or channel).',
		required: false,
	}),
	user: flags.string({
		char: 'u',
		description: 'The unique identifier of the user to ban.',
		required: false,
	}),
	reason: flags.string({
		char: 'r',
		description: 'A reason for adding a timeout.',
		required: false,
	}),
	duration: flags.string({
		char: 'd',
		description: 'Duration of timeout in minutes.',
		default: '60',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.UserBan = UserBan;
