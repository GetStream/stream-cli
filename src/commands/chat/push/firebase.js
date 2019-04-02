const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class PushFirebase extends Command {
	async run() {
		const { flags } = this.parse(PushFirebase);

		try {
			if (!flags.disable && !flags.key) {
				const res = await prompt([
					{
						type: 'input',
						name: 'key',
						message: `What is your API key for Firebase?`,
						required: true,
					},
					{
						type: 'input',
						name: 'notification_template',
						hint: 'Omit for Stream default',
						message: `What JSON notification template would you like to use?`,
						required: false,
					},
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await chatAuth(this);

			const payload = {
				firebase_config: {
					api_key: flags.key,
					disabled: flags.disable || false,
				},
			};

			if (flags.notification_template) {
				payload.firebase_config.notification_template =
					flags.notification_template;
			}

			await client.updateAppSettings(payload);

			if (flags.json) {
				const settings = await client.getAppSettings();

				this.log(
					JSON.stringify(settings.app.push_notifications.firebase)
				);
				this.exit();
			}

			this.log(
				`Push notifications have been ${
					flags.disable
						? chalk.red('disabled')
						: chalk.green('enabled')
				} for ${chalk.bold('Firebase')}.`
			);
			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

PushFirebase.flags = {
	key: flags.string({
		char: 'k',
		description: 'API key for Firebase.',
		required: false,
	}),
	notification_template: flags.string({
		char: 'n',
		description: 'JSON notification template.',
		required: false,
	}),
	disable: flags.boolean({
		description: 'Disable Firebase push notifications and clear config.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.PushFirebase = PushFirebase;
