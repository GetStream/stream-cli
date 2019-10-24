const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class ChannelMarkRead extends Command {
	async run() {
		const { flags } = this.parse(ChannelMarkRead);

		try {
			if (!flags.user || !flags.channel) {
				const res = await prompt([
					{
						type: 'input',
						name: 'user',
						message: `What is the unique identifier for the user marking messages as read?`,
						required: true,
					},
					{
						type: 'input',
						name: 'channel',
						message: `What is the unique identifier for the channel?`,
						required: true,
					},
				]);

				flags.channel = res.channel;
			}

			const client = await chatAuth(this);
			const channel = await client.channel(flags.type, flags.channel);

			await channel.markAllRead({ user: { id: flags.user } });

			this.log(
				`Channel ${chalk.bold(
					flags.channel
				)} messages have been marked as read.`
			);
			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

ChannelMarkRead.flags = {
	user: flags.string({
		char: 'u',
		description: 'The ID of the user marking all messages as read.',
		required: false,
	}),
	channel: flags.string({
		char: 'c',
		description: 'The ID of the channel you wish to update.',
		required: false,
	}),
	type: flags.string({
		char: 't',
		description: 'Type of channel.',
		required: false,
	}),
};

module.exports.ChannelMarkRead = ChannelMarkRead;
