const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class MessageRemove extends Command {
	async run() {
		const { flags } = this.parse(MessageRemove);

		try {
			if (!flags.message) {
				const res = await prompt([
					{
						type: 'input',
						name: 'message',
						message: `What is the unique identifier for the message?`,
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
			const remove = await client.deleteMessage(flags.message);

			if (flags.json) {
				this.log(JSON.stringify(remove));
				this.exit();
			}

			this.log(
				`The message ${chalk.bold(flags.message)} has been removed.`
			);
			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

MessageRemove.flags = {
	message: flags.string({
		char: 'message',
		description:
			'The unique identifier of the message you would like to remove.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.MessageRemove = MessageRemove;
