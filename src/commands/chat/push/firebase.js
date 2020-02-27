import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';

class PushFirebase extends Command {
	async run() {
		const { flags } = this.parse(PushFirebase);

		try {
			const client = await chatAuth(this);
			if (flags.disable) {
				const result = await prompt({
					type: 'toggle',
					name: 'proceed',
					message:
						'This will disable Firebase push notifications and remove your Firebase Server Key. Are you sure?',
					required: true
				});
				if (result.proceed) {
					await client.updateAppSettings({
						firebase_config: {
							disabled: true
						}
					});
					this.log(`Push notifications have been ${chalk.red('disabled')} with ${chalk.bold('Firebase')}.`);
				}
				this.exit();
			} else if (!flags.key) {
				const res = await prompt([
					{
						type: 'input',
						name: 'key',
						message: `What is your Server key for Firebase?`,
						required: true
					},
					{
						type: 'input',
						name: 'notification_template',
						hint: 'Omit for Stream default',
						message: `What JSON notification template would you like to use?`,
						required: false
					},
					{
						type: 'input',
						name: 'data_template',
						hint: 'Omit for Stream default',
						message: `What JSON data template would you like to use?`,
						required: false
					}
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const payload = {
				firebase_config: {
					server_key: flags.key
				}
			};

			if (flags.notification_template) {
				payload.firebase_config.notification_template = flags.notification_template;
			}
			if (flags.data_template) {
				payload.firebase_config.data_template = flags.data_template;
			}

			await client.updateAppSettings(payload);

			if (flags.json) {
				const settings = await client.getAppSettings();

				this.log(JSON.stringify(settings.app.push_notifications.firebase));
				this.exit();
			}

			this.log(`Push notifications have been ${chalk.green('enabled')} for ${chalk.bold('Firebase')}.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

PushFirebase.flags = {
	key: flags.string({
		char: 'k',
		description: 'Server key for Firebase.',
		required: false
	}),
	notification_template: flags.string({
		char: 'n',
		description: 'JSON notification template.',
		required: false
	}),
	data_template: flags.string({
		char: 'd',
		description: 'JSON data template.',
		required: false
	}),
	disable: flags.boolean({
		description: 'Disable Firebase push notifications and clear config.',
		required: false
	}),
	json: flags.boolean({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

PushFirebase.description = 'Specifies Firebase for push notifications.';

module.exports.PushFirebase = PushFirebase;
