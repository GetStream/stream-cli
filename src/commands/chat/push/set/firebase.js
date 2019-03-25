const { Command, flags } = require('@oclif/command');

const { auth } = require('../../../../utils/auth');

class PushSetFirebase extends Command {
	async run() {
		const { flags } = this.parse(PushSetFirebase);

		try {
			const client = await auth(this);

			const settings = await client.updateAppSettings({
				firebase_config: {
					api_key: flags.api_key,
					notification_template: flags.notification_template,
				},
			});

			if (flags.json) {
				this.log(JSON.stringify(settings));
				this.exit();
			}

			this.log('Push notifications have been enabled for Firebase.');
			this.exit(0);
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

PushSetFirebase.flags = {
	api_key: flags.string({
		char: 'f',
		description: 'API key for Firebase.',
		required: false,
	}),
	notification_template: flags.string({
		char: 'n',
		description: 'JSON template for notifications (APN and Firebase).',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.SettingsPush = PushSetFirebase;
