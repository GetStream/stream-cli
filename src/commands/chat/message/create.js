const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');
const uuid = require('uuid/v4');

const { chatAuth } = require('../../../utils/auth/chat-auth');
const { credentials } = require('../../../utils/config');

class MessageCreate extends Command {
	async run() {
		const { flags } = this.parse(MessageCreate);

		try {
			const { name } = await credentials(this);

			if (
				!flags.user ||
				!flags.name ||
				!flags.channel ||
				!flags.type ||
				!flags.message
			) {
				const res = await prompt([
					{
						type: 'input',
						name: 'user',
						message: `What is the unique identifier for the user sending this message?`,
						default: uuid(),
						required: true,
					},
					{
						type: 'input',
						name: 'name',
						message: `What is the name of the user sending this message?`,
						default: name,
						required: true,
					},
					{
						type: 'input',
						name: 'image',
						message: `What is an absolute URL to the avatar of the user sending this message?`,
						required: false,
					},
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
						name: 'message',
						message: 'What is the message you would like to send?',
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

			const create = await channel.sendMessage({
				text: flags.message,
				user: {
					id: flags.user,
					name: flags.name,
					image: flags.image || null,
				},
			});

			if (flags.json) {
				this.log(JSON.stringify(create.message));
				this.exit();
			}

			this.log(`Message ${chalk.bold(create.message.id)} was created.`);
			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

MessageCreate.flags = {
	user: flags.string({
		char: 'u',
		description: 'The ID of the user sending the message.',
		required: false,
	}),
	name: flags.string({
		char: 'n',
		description: 'The name of the user sending the message.',
		required: false,
	}),
	image: flags.string({
		char: 'i',
		description:
			'Absolute URL for an avatar of the user sending the message.',
		required: false,
	}),
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
	message: flags.string({
		char: 'm',
		description: 'The message you would like to send as plaintext.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.MessageCreate = MessageCreate;
