const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class ChannelDemoteModerator extends Command {
	async run() {
		const { flags } = this.parse(ChannelDemoteModerator);

		try {
			if (!flags.channel || !flags.type || !flags.image) {
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
						message: `What is the unique ID of the user to promote as a moderator?`,
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
			const channel = await client.channel(flags.type, flags.channel);

			const exists = await client.queryUsers({
				id: { $in: [flags.user] },
			});

			if (!exists.users.length) {
				this.log(
					`The user ${flags.user} in channel ${chalk.bold(
						flags.channel
					)} (${flags.type}) does not exist.`
				);

				this.exit();
			}

			await channel.addModerators([flags.user]);

			if (flags.json) {
				this.log(JSON.stringify(create.channel));
				this.exit();
			}

			this.log(`User ${chalk.bold(flags.user)} has been promoted.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

ChannelDemoteModerator.flags = {
	channel: flags.string({
		char: 'c',
		description: 'A unique ID for the channel you wish to create.',
		required: false,
	}),
	type: flags.string({
		char: 't',
		description: 'Type of channel.',
		required: false,
	}),
	user: flags.string({
		char: 'u',
		description: 'A unique ID for user user to demote.',
		required: false,
	}),
};

module.exports.ChannelDemoteModerator = ChannelDemoteModerator;
