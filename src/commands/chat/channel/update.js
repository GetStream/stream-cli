const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { auth } = require('../../../utils/auth');
const { credentials } = require('../../../utils/config');

class ChannelUpdate extends Command {
	async run() {
		const { flags } = this.parse(ChannelUpdate);

		try {
			const { name } = await credentials(this);

			if (!flags.channel || !flags.type) {
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
						name: 'name',
						message: `What is the new name for the channel?`,
						required: false,
					},
					{
						type: 'input',
						name: 'image',
						message: `What is the absolute image URL for the channel?`,
						required: false,
					},
					{
						type: 'input',
						name: 'description',
						message: `What description would you like to set for the channel?`,
						required: false,
					},
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await auth(this);
			const channel = await client.channel(flags.type, flags.channel);

			const payload = {
				updated_by: {
					id: 'CLI',
					name,
				},
			};

			if (flags.name) payload.name = flags.name;
			if (flags.image) payload.image = flags.image;
			if (flags.description) payload.description = flags.description;

			const update = await channel.update(payload);

			if (flags.json) {
				this.log(JSON.stringify(update));
				this.exit(0);
			}

			this.log(`Channel ${chalk.bold(flags.channel)} has been updated.`);
			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

ChannelUpdate.flags = {
	channel: flags.string({
		char: 'c',
		description: 'The ID of the channel you wish to update.',
		required: false,
	}),
	type: flags.string({
		char: 't',
		description: 'Type of channel.',
		options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
		required: false,
	}),
	name: flags.string({
		char: 'n',
		description: 'Name of the channel room.',
		required: false,
	}),
	image: flags.string({
		char: 'i',
		description: 'URL to the channel image.',
		required: false,
	}),
	description: flags.string({
		char: 'd',
		description: 'Description for the channel.',
		required: false,
	}),
	reason: flags.string({
		char: 'r',
		description: 'Reason for changing channel.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.ChannelUpdate = ChannelUpdate;
