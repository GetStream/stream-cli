const { Command, flags } = require('@oclif/command');
const fs = require('fs');

const { auth } = require('../../../utils/auth');

class SettingsPush extends Command {
	async run() {
		const { flags } = this.parse(SettingsPush);

		try {
			if (!flags.type) {
				this.log(
					'Please specify a push notification type of apn, firebase or webhook.'
				);
				this.exit(0);
			}

			const client = await auth(this);

			if (flags.enable && flags.type === 'apn') {
				await client.updateAppSettings({
					apn_config: {
						p12_cert: fs.readFileSync(flags.p12_cert) || '',
						pem_cert:
							fs.readFileSync(flags.pem_cert, 'utf-8') || '',
						auth_key: flags.auth_key || '',
						key_id: flags.key_id || '',
						team_id: flags.team_id || '',
						topic: flags.bundle_id || '',
						notification_template: flags.notification_template,
					},
				});

				if (!flags.json) {
					this.log('Push notifications have been enabled with APN.');
				}
			}

			if (flags.enable && flags.type === 'firebase') {
				await client.updateAppSettings({
					firebase_config: {
						api_key: flags.api_key,
						notification_template: flags.notification_template,
					},
				});

				if (!flags.json) {
					this.log(
						'Push notifications have been enabled for Firebase.'
					);
				}
			}

			if (flags.enable && flags.type === 'webhook') {
				await client.updateAppSettings({
					webhook_url: flags.webhook_url,
				});

				if (!flags.json) {
					this.log(
						'Push notifications have been enabled for Webhooks.'
					);
				}
			}

			if (flags.disable) {
				this.log(`Push via ${flags.type} has been disabled.`);
			}

			this.exit(0);
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

SettingsPush.flags = {
	enable: flags.boolean({
		chart: 'e',
		description: 'Enable push notifications for your project.',
		required: false,
	}),
	disable: flags.boolean({
		chart: 'd',
		description: 'Disable push notifications for your project.',
		required: false,
	}),
	type: flags.string({
		char: 't',
		description: 'Type of configuration.',
		options: ['apn', 'firebase', 'webhook'],
		required: false,
	}),
	auth_key: flags.string({
		char: 'a',
		description: 'Private auth key for APN.',
		required: false,
	}),
	key_id: flags.string({
		char: 'k',
		description: 'Key ID for APN.',
		required: false,
	}),
	bundle_id: flags.string({
		char: 'b',
		description: 'Bundle identifier for APN (e.g. com.apple.test)',
		required: false,
	}),
	team_id: flags.string({
		char: 'i',
		description: 'Team ID for APN.',
		required: false,
	}),
	pem_cert: flags.string({
		char: 'p',
		description: 'Absolute path to RSA key for APN (.pem).',
		required: false,
	}),
	p12_cert: flags.string({
		char: 'b',
		description: 'Absolute path to .p12 file for APN.',
		required: false,
	}),
	notification_template: flags.string({
		char: 'n',
		description: 'JSON template for notifications (APN and Firebase).',
		required: false,
	}),
	api_key: flags.string({
		char: 'f',
		description: 'API key for Firebase.',
		required: false,
	}),
	webhook_url: flags.string({
		char: 'w',
		description: 'Fully qualified URL for webhook support.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.SettingsPush = SettingsPush;
