import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';

class ChannelHide extends Command {
	async run() {
		const { flags } = this.parse(ChannelHide);

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

			const channel = client.channel(flags.type, flags.channel);
			await channel.hide();

			this.log(`The channel ${chalk.bold(flags.channel)} has been hidden.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

ChannelHide.flags = {
	channel: flags.string({
		char: 'c',
		description: 'The channel ID you wish to remove.',
		required: false
	}),
	type: flags.string({
		char: 't',
		description: 'Type of channel.',
		required: false
	})
};

ChannelHide.description = 'Hides a channel.';

module.exports.ChannelHide = ChannelHide;
