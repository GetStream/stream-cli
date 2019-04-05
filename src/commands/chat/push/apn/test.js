const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const Table = require('cli-table');
const chalk = require('chalk');

const { chatAuth } = require('../../../../utils/auth/chat-auth');

class APNPushTest extends Command {
	async run() {
		const { flags } = this.parse(APNPushTest);

		try {
			if (!flags.user_id) {
				const result = await prompt([
					{
						type: 'input',
						name: 'user_id',
						hint: 'user-123',
						message:
							'What is the ID of the user you wish to test push settings with?',
						required: true,
					},
					{
						type: 'input',
						name: 'message_id',
						hint: `Omit to pick random message`,
						message:
							'What is the message ID you want to use to render the notification template with?',
						required: false,
					},
					{
						type: 'input',
						name: 'notification_template',
						hint: 'Omit for the template configured in your app',
						message: `What JSON notification template would you like to use?`,
						required: false,
					},
				]);

				for (const key in result) {
					if (result.hasOwnProperty(key)) {
						flags[key] = result[key];
					}
				}
			}

			const client = await chatAuth(this);

			const payload = {
				messageID: flags.message_id || '',
				apnTemplate: flags.notification_template || '',
			};
			const userID = flags.user_id || '';

			const response = await client.testPushSettings(userID, payload);

			if (flags.json) {
				this.log(JSON.stringify(response));
				this.exit();
			}

			if (response.general_errors) {
				this.log(
					`It seems there were some ${chalk.red.bold(
						'errors'
					)} with the input you provided. Listing them below:`
				);
				for (const err of response.general_errors) {
					this.log(`\t - ${err}`);
				}
			}

			if (response.rendered_apn_template) {
				this.log(
					`Here is the rendered notification that was sent to your devices:`
				);
				this.log(
					JSON.stringify(
						JSON.parse(response.rendered_apn_template),
						null,
						4
					)
				);
			}

			if (response.device_errors) {
				this.log(
					`It seems we couldn't push the notification to all of the user devices. Here are the ${chalk.red.bold(
						'errors'
					)} for each device:`
				);
				const table = new Table();
				for (const [deviceID, deviceError] of Object.entries(
					response.device_errors
				)) {
					table.push({ [deviceID]: deviceError });
				}
				this.log(table.toString());
			}

			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

APNPushTest.flags = {
	user_id: flags.string({
		char: 'u',
		description: 'User ID',
		required: false,
	}),
	message_id: flags.string({
		char: 'm',
		description: 'Message ID.',
		required: false,
	}),
	notification_template: flags.string({
		char: 'n',
		description: 'Notification template',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.PushTest = APNPushTest;
