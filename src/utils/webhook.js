import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import chalk from 'chalk';

import { chatAuth } from 'utils/auth/chat-auth';




module.exports.setFlags = (command, description) => {
	command.flags = {
		url: flags.string({
			char: 'u',
			description: 'A fully qualified URL for webhook support.',
			required: false
		}),
		json: flags.boolean({
			char: 'j',
			description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
			required: false
		})
	};
	command.description = description
}

module.exports.runWebhookCommand = async (appFieldName, command, flags) => {
	try {
		if (!flags.url) {
			const res = await prompt([
				{
					type: 'input',
					name: 'url',
					message: `What is the absolute URL for your webhook?`,
					required: true
				}
			]);

			flags.url = res.url;
		}

		let settings = {
			[appFieldName]: flags.url
		}

		const client = await chatAuth(command);
		await client.updateAppSettings(settings);

		if (flags.json) {
			const settings = await client.getAppSettings();

			command.log(
				JSON.stringify(settings)
			);
			command.exit();
		}

		command.log(`Webhook ${chalk.bold(appFieldName)} is set to ${flags.url}`);
		command.exit();
	} catch (error) {
		await command.config.runHook('telemetry', {
			ctx: command,
			error
		});
	}
}
