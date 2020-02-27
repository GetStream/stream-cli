import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';
import { v4 as uuid } from 'uuid';

import { chatAuth } from 'utils/auth/chat-auth';
import { credentials } from 'utils/config';

class ChannelCreate extends Command {
	async run() {
		const { flags } = this.parse(ChannelCreate);

		try {
			if (!flags.channel || !flags.type) {
				const res = await prompt([
					{
						type: 'input',
						name: 'channel',
						message: `What is the unique identifier for the channel?`,
						default: uuid(),
						required: true
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
							{ message: 'Team', value: 'team' }
						]
					},
					{
						type: 'input',
						name: 'name',
						message: `What is the name of your channel?`,
						default: uuid(),
						required: false
					},
					{
						type: 'input',
						name: 'image',
						message: `What is the absolute URL to the channel image?`,
						hint: 'optional',
						required: false
					}
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const { name } = await credentials(this);
			const client = await chatAuth(this);

			let payload = {
				name: flags.name,
				created_by: {
					id: 'CLI',
					name
				}
			};
			if (flags.image) payload.image = flags.image;

			if (flags.data) {
				const parsed = JSON.parse(flags.data);
				payload = Object.assign({}, payload, parsed);
			}

			const channel = await client.channel(flags.type, flags.channel, payload);

			const create = await channel.create();

			if (flags.json) {
				this.log(JSON.stringify(create.channel));
				this.exit();
			}

			this.log(`Channel ${chalk.bold(create.channel.id)} has been created.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

ChannelCreate.flags = {
	channel: flags.string({
		char: 'c',
		description: 'A unique ID for the channel you wish to create.',
		default: uuid(),
		required: false
	}),
	type: flags.string({
		char: 't',
		description: 'Type of channel.',
		required: false
	}),
	name: flags.string({
		char: 'n',
		description: 'Name of the channel room.',
		required: false
	}),
	image: flags.string({
		char: 'i',
		description: 'URL to channel image.',
		required: false
	}),
	users: flags.string({
		char: 'u',
		description: 'Comma separated list of users to add.',
		required: false
	}),
	data: flags.string({
		char: 'd',
		description: 'Additional data as JSON.',
		required: false
	}),
	json: flags.boolean({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

ChannelCreate.description = 'Creates a new channel.';

module.exports.ChannelCreate = ChannelCreate;
