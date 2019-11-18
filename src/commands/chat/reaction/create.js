const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class ReactionCreate extends Command {
	async run() {
		const { flags } = this.parse(ReactionCreate);

		try {
			if (!flags.json) {
				const res = await prompt([
					{
						type: 'input',
						name: 'channel',
						message:
							'What is the unique identifier for the channel?',
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
						message:
							'What is the unique identifier for the message?',
						required: true,
					},
					{
						type: 'input',
						name: 'reaction',
						hint: 'love',
						message: 'What is the reaction you would like to add?',
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

			const channel = client.channel(flags.type, flags.channel);
			const reaction = await channel.sendReaction(flags.message, {
				type: flags.reaction,
			});

			if (this.json) {
				this.log(JSON.stringify(reaction));
				this.exit();
			}

			this.log('Your reaction has been created.');
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

ReactionCreate.flags = {
	channel: flags.string({
		char: 'c',
		description: 'The unique identifier for the channel.',
		required: false,
	}),
	type: flags.string({
		char: 't',
		description: 'The type of channel.',
		required: false,
	}),
	message: flags.string({
		char: 'c',
		description: 'The unique identifier for the message.',
		required: false,
	}),
	reaction: flags.string({
		char: 'r',
		description: 'A reaction for the message (e.g. love).',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

ReactionCreate.description = 'Creates a new reaction.';

module.exports.ReactionCreate = ReactionCreate;
