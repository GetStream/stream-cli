import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';

class ChannelGet extends Command {
	async run() {
		const { flags } = this.parse(ChannelGet);

		try {
			if (!flags.channel || !flags.type) {
				const res = await prompt([
					{
						type: 'input',
						name: 'channel',
						message: `What is the unique identifier for the channel?`,
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
					}
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await chatAuth(this);

			const channel = await client.queryChannels(
				{ id: flags.channel, type: flags.type },
				{ last_message_at: -1 },
				{
					subscribe: false
				}
			);

			if (!channel.length) {
				this.log(
					`Channel ${chalk.bold(flags.channel)} with type ${chalk.bold(flags.type)} could not be found.`
				);
				this.exit();
			}

			this.log(JSON.stringify(channel[0].data));
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

ChannelGet.flags = {
	channel: flags.string({
		char: 'c',
		description: 'The channel ID you wish to retrieve.',
		required: false
	}),
	type: flags.string({
		char: 't',
		description: 'Type of channel.',
		required: false
	})
};

ChannelGet.description = 'Gets a specific channel by its ID and type.';

module.exports.ChannelGet = ChannelGet;
