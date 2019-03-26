const { Command, flags } = require('@oclif/command');
const Table = require('cli-table');
const chalk = require('chalk');

const { auth } = require('../../../utils/auth');

class PushGet extends Command {
	async run() {
		const { flags } = this.parse(PushGet);

		try {
			const client = await auth(this);

			const settings = await client.getAppSettings();

			if (flags.json) {
				this.log(JSON.stringify(settings.app));
				this.exit();
			}

			const table = new Table();

			table.push(
				{
					[`${chalk.green.bold('APN - Status')}`]:
						settings.app.push_notifications.apn.enabled ===
						'enabled'
							? 'Enabled'
							: 'Disabled',
				},
				{
					[`${chalk.green.bold('APN – Auth Type')}`]: !settings.app
						.push_notifications.apn.auth_type
						? 'N/A'
						: settings.app.push_notifications.apn.auth_type.toUpperCase(),
				},
				{
					[`${chalk.green.bold(
						'APN – Notification Template'
					)}`]: !settings.app.push_notifications.apn
						.notification_template
						? 'Stream Default'
						: settings.app.push_notifications.apn
								.notification_template,
				},
				{
					[`${chalk.green.bold('Firebase')}`]:
						settings.app.push_notifications.firebase.enabled ===
						'enabled'
							? 'Enabled'
							: 'Disabled',
				},
				{
					[`${chalk.green.bold(
						'Firebase – Notification Template'
					)}`]: !settings.app.push_notifications.firebase
						.notification_template
						? 'Stream Default'
						: settings.app.push_notifications.firebase
								.notification_template,
				},
				{
					[`${chalk.green.bold('Webhook – URL')}`]: !settings.app
						.webhook_url
						? 'N/A'
						: settings.app.webhook_url,
				}
			);

			this.log(table.toString());
			this.exit(0);
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

PushGet.flags = {
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.SettingsGet = PushGet;
