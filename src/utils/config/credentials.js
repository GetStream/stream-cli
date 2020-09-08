const emoji = require('node-emoji');
const chalk = require('chalk');
import path from 'path';
import fs from 'fs-extra';

async function credentials(ctx) {
	const config = path.join(ctx.config.configDir, 'config.json');

	try {
		if (!(await fs.pathExists(config))) {
			await fs.outputJson(config, {
				name: '',
				email: '',
				apiKey: '',
				apiSecret: '',
				apiBaseUrl: 'https://chat-us-east-1.stream-io-api.com',
				environment: 'production',
				telemetry: true,
				timeout: 3000,
			});
		}

		const {
			name,
			email,
			apiKey,
			apiSecret,
			apiBaseUrl,
			environment,
			telemetry,
			timeout,
		} = await fs.readJson(config);

		if (
			!name ||
			!email ||
			!apiKey ||
			!apiSecret ||
			!apiBaseUrl ||
			!environment ||
			!telemetry ||
			!timeout
		) {
			console.warn(
				`Credentials not found. Run the command ${chalk.bold(
					'stream config:set'
				)} to generate a new Stream configuration file.`,
				emoji.get('warning')
			);

			ctx.exit(0);
		}

		return {
			name,
			email,
			apiKey,
			apiSecret,
			apiBaseUrl,
			environment,
			telemetry,
			timeout,
		};
	} catch (error) {
		ctx.error(error || 'A Stream CLI error has occurred.', {
			exit: 1,
		});
	}
}

module.exports.credentials = credentials;
