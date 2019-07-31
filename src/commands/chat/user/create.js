const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class UserCreate extends Command {
	async run() {
		const { flags } = this.parse(UserCreate);

		try {
			if (!flags.type || !flags.channel || !flags.user || !flags.role) {
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
						name: 'user',
						message: `What is the unique identifier for the user?`,
						required: true,
					},
					{
						type: 'select',
						name: 'role',
						message: 'What role would you like assign to the user?',
						required: true,
						choices: [
							{
								message: 'User',
								value: 'user',
							},
							{
								message: 'Moderator',
								value: 'moderator',
							},
							{
								message: 'Guest',
								value: 'guest',
							},
							{
								message: 'Admin',
								value: 'admin',
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
			await client.updateUser({ id: flags.user, role: flags.role });

			const create = await client
				.channel(flags.type, flags.channel)
				.addMembers([flags.user]);

			if (flags.json) {
				this.log(JSON.stringify(create.channel.members[0].user));
				this.exit();
			}

			this.log(create.channel.members[0]);
			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

UserCreate.flags = {
	channel: flags.string({
		char: 'c',
		description: 'Channel identifier.',
		required: false,
	}),
	type: flags.string({
		char: 't',
		description: 'The type of channel.',
		required: false,
	}),
	user: flags.string({
		char: 'u',
		description: 'Comma separated list of users to add.',
		required: false,
	}),
	role: flags.string({
		char: 'r',
		description: 'The role to assign to the user.',
		options: [
			'admin',
			'guest',
			'channel_moderator',
			'channel_member',
			'channel_owner',
			'message_owner',
		],
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.UserCreate = UserCreate;
