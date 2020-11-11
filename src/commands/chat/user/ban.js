import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';

class UserBan extends Command {
	async run() {
		const { flags } = this.parse(UserBan);

		try {
			let type = flags.type;
			if (!type) {
				let res = await prompt([
					{
						type: 'select',
						name: 'type',
						message: 'Would you like to apply a global or channel ban?',
						required: true,
						choices: [ { message: 'Global', value: 'global' }, { message: 'Channel', value: 'channel' } ]
					}
				]);
				type = res.type;
			}

			let cid = flags.cid;
			if (type === 'channel' && !cid) {
				cid = await prompt([
					{
						type: 'input',
						name: 'id',
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
			}

			if (!flags.user || !flags.reason || !flags.duration) {
				const res = await prompt([
					{
						type: 'input',
						name: 'user',
						message: 'What is the unique identifier for the user?',
						required: true
					},
					{
						type: 'input',
						name: 'reason',
						message: 'What is the reason for banning the user?',
						required: true
					},
					{
						type: 'input',
						name: 'duration',
						hint: 'minutes',
						message: 'How long would you like to ban the user for?',
						required: false
					}
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			if (!flags.reason.length) {
				this.error('A reason is required to ban a user.');
				this.exit();
			}

			const client = await chatAuth(this);
			const payload = {
				reason: flags.reason,
				user_id: 'CLI'
			};

			await client.updateUser({
				id: 'CLI',
				name: 'CLI',
				role: 'admin'
			});

			if (flags.duration) {
				payload.timeout = parseInt(flags.duration * 60, 10);
			}
			payload.ip_ban = flags.ip;

			let ban;
			if (type === 'channel') {
				ban = await client.channel(cid.type, cid.id).banUser(flags.user, payload);
			} else if (type === 'global') {
				this.log(`Global ban`);
				ban = await client.banUser(flags.user, payload);
			} else {
				this.warn("invalid ban type")
			}

			if (flags.json) {
				this.log(JSON.stringify(ban));
				this.exit();
			}

			this.log(`The user ${chalk.bold(flags.user)} has been banned.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

UserBan.flags = {
	type: flags.string({
		char: 't',
		description: 'Type of ban to perform (e.g. global or channel).',
		required: false
	}),
	user: flags.string({
		char: 'u',
		description: 'The unique identifier of the user to ban.',
		required: false
	}),
	reason: flags.string({
		char: 'r',
		description: 'A reason for adding a timeout.',
		required: false
	}),
	duration: flags.string({
		char: 'd',
		description: 'Duration of timeout in minutes.',
		default: '60',
		required: false
	}),
	ip: flags.boolean({
		description: 'Apply IP ban as well',
		default: false,
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

UserBan.description = 'Bans a user.';

module.exports.UserBan = UserBan;
