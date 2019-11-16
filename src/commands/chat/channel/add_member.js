const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');
const { credentials } = require('../../../utils/config');

class ChannelAddMember extends Command {
	async run() {
		const { flags } = this.parse(ChannelAddMember);

		try {
			if (!flags.channel || !flags.type || !flags.image) {
				const res = await prompt([
					{
						type: 'input',
						name: 'channel',
						message: 'What is the unique ID for the channel?',
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
						message: 'What is the unique ID of the user to add?',
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

			if (typeof channel === 'object') {
				await channel.addMembers([flags.user]);

				if (flags.json) {
					this.log(JSON.stringify(channel));
					this.exit();
				}

				this.log(
					`User ${chalk.bold(flags.user)} has been added as a member.`
				);

				this.exit();
			} else if (!Array.isArray(channel) && !channel.length) {
				this.log(
					`Channel ${chalk.bold(
						flags.channel
					)} with type ${chalk.bold(flags.type)} could not be found.`
				);

				this.exit();
			}
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

ChannelAddMember.flags = {
	channel: flags.string({
		char: 'c',
		description: 'A unique ID for the channel add the user to.',
		required: false,
	}),
	type: flags.string({
		char: 't',
		description: 'Type of channel.',
		required: false,
	}),
	name: flags.string({
		char: 'n',
		description: 'Name of the channel room.',
		required: false,
	}),
	image: flags.string({
		char: 'i',
		description: 'URL to channel image.',
		required: false,
	}),
	users: flags.string({
		char: 'u',
		description: 'Unique identifier for the user you are adding.',
		required: false,
	}),
	data: flags.string({
		char: 'r',
		description: 'The role of the user you are adding.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.ChannelAddMember = ChannelAddMember;
