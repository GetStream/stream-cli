const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class PushWebhook extends Command {
	async run() {
		const { flags } = this.parse(PushWebhook);

		try {
			if (!flags.url) {
				const res = await prompt([
					{
						type: 'input',
						name: 'url',
						message: `What is the absolute URL for your webhook?`,
						required: true,
					},
				]);

				flags.url = res.url;
			}

			const client = await chatAuth(this);
			await client.updateAppSettings({
				webhook_url: flags.url,
			});

			if (flags.json) {
				const settings = await client.getAppSettings();

				this.log(
					JSON.stringify({
						webhook_url: settings.app.webhook_url,
					})
				);
				this.exit();
			}

			this.log(
				`Push notifications have been enabled for ${chalk.bold(
					'Webhooks'
				)}.`
			);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

PushWebhook.flags = {
	url: flags.string({
		char: 'u',
		description: 'A fully qualified URL for webhook support.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

PushWebhook.description = 'Tests webhook notifications.';

module.exports.PushWebhook = PushWebhook;
