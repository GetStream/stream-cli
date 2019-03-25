const { Command, flags } = require('@oclif/command');

const { auth } = require('../../../../utils/auth');

class PushSetWebhook extends Command {
	async run() {
		const { flags } = this.parse(PushSetWebhook);

		try {
			const client = await auth(this);

			const settings = await client.updateAppSettings({
				webhook_url: flags.webhook_url,
			});

			if (flags.json) {
				this.log(JSON.stringify(settings));
			}

			this.log('Push notifications have been enabled for Webhooks.');
			this.exit(0);
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

PushSetWebhook.flags = {
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

module.exports.SettingsPush = PushSetWebhook;
