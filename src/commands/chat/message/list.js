const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const moment = require('moment');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class MessageList extends Command {
	async run() {
		const { flags } = this.parse(MessageList);

		try {
			if (!flags.channel || !flags.type || !flags.json) {
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
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await chatAuth(this);
			client.channel(flags.type, flags.channel);

			const messages = await client.queryChannels(
				{},
				{},
				{ watch: false, presence: false }
			);

			if (flags.json) {
				if (messages.length === 0) {
					this.log(JSON.stringify(messages));
					this.exit();
				}

				this.log(JSON.stringify(messages[0].state.messages));
				this.exit();
			}

			const data = messages[0].state.messages;

			if (data.length === 0) {
				this.log('No messages available.');
				this.exit();
			}

			for (let i = 0; i < data.length; i++) {
				const timestamp = `${
					data[i].deleted_at ? 'Deleted on' : 'Created at'
				} ${moment(data[i].created_at).format(
					'dddd, MMMM Do YYYY [at] h:mm:ss A'
				)}`;

				this.log(
					`Message ${chalk.bold(data[i].id)} (${timestamp}): ${
						data[i].text
					}`
				);
			}

			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

MessageList.flags = {
	type: flags.string({
		char: 't',
		description: 'The type of channel.',
		required: false,
	}),
	channel: flags.string({
		char: 'c',
		description:
			'The ID of the channel that you would like to send a message to.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.MessageList = MessageList;
