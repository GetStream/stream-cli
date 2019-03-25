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
				this.log(JSON.stringify(settings.app.push_notifications));
				this.exit();
			}

			const push = settings.app.push_notifications;

			const table = new Table();

			table.push(
				{
					[`${chalk.green.bold('APN - Status')}`]:
						push.apn.enabled === 'enabled' ? 'Enabled' : 'Disabled',
				},
				{
					[`${chalk.green.bold('APN – Auth Type')}`]: !push.apn
						.auth_type
						? 'N/A'
						: push.apn.auth_type.toUpperCase(),
				},
				{
					[`${chalk.green.bold(
						'APN – Notification Template'
					)}`]: !push.apn.notification_template
						? 'Stream Default'
						: push.apn.notification_template,
				},
				{
					[`${chalk.green.bold('Firebase')}`]:
						push.firebase.enabled === 'enabled'
							? 'Enabled'
							: 'Disabled',
				},
				{
					[`${chalk.green.bold('Firebase – Notification Template')}`]:
						push.firebase.notification_template === ''
							? 'Stream Default'
							: push.firebase.notification_template,
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
