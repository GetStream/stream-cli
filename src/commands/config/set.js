import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import emoji from 'node-emoji';
import path from 'path';
import fs from 'fs-extra';

import { credentials } from '../../utils/config';

class ConfigSet extends Command {
	async run() {
		const { flags } = this.parse(ConfigSet);
		const config = path.join(this.config.configDir, 'config.json');

		try {
			if (!flags.name || !flags.email || !flags.key || !flags.secret || !flags.url || !flags.environment) {
				const res = await prompt([
					{
						type: 'input',
						name: 'name',
						message: `What is your full name?`,
						required: true
					},
					{
						type: 'input',
						name: 'email',
						message: `What is your email address associated with Stream?`,
						required: true
					},
					{
						type: 'input',
						name: 'key',
						message: `What is your Stream API key?`,
						required: true
					},
					{
						type: 'password',
						name: 'secret',
						message: `What is your Stream API secret?`,
						required: true
					},
					{
						type: 'input',
						name: 'url',
						message: `What is your Stream API base URL?`,
						default: 'https://chat-us-east-1.stream-io-api.com',
						required: false
					},
					{
						type: 'input',
						name: 'environment',
						message: `What environment would you like to run in?`,
						default: 'production',
						required: false
					},
					{
						type: 'input',
						name: 'telemetry',
						message: `Would you like to enable error tracking for debugging purposes?`,
						default: true,
						required: false
					}
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			await fs.ensureDir(this.config.configDir);
			await fs.writeJson(config, {
				name: flags.name,
				email: flags.email.toLowerCase(),
				apiKey: flags.key,
				apiSecret: flags.secret,
				apiBaseUrl: flags.url,
				environment: flags.environment,
				telemetry: flags.telemetry
			});

			if (flags.json) {
				this.log(JSON.stringify(await credentials(this)));
				this.exit();
			}

			this.log('Your Stream CLI configuration has been generated!', emoji.get('rocket'));
			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1
			});
		}
	}
}

ConfigSet.flags = {
	name: flags.string({
		char: 'n',
		description: 'Full name for configuration.',
		required: false
	}),
	email: flags.string({
		char: 'e',
		description: 'Email for configuration.',
		required: false
	}),
	key: flags.string({
		char: 'k',
		description: 'API key for configuration.',
		required: false
	}),
	secret: flags.string({
		char: 's',
		description: 'API secret for configuration.',
		required: false
	}),
	url: flags.string({
		char: 'u',
		description: 'API base URL for configuration.',
		required: false
	}),
	environment: flags.string({
		char: 'v',
		description: 'Environment to run in (production or development for token and permission checking).',
		required: false
	}),
	telemetry: flags.boolean({
		char: 't',
		description: 'Enable error reporting for debugging purposes.',
		required: false
	}),
	json: flags.boolean({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

ConfigSet.description = 'Sets your user configuration.';

module.exports.ConfigSet = ConfigSet;
