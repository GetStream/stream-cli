import { Command, flags } from '@oclif/command';
import Table from 'cli-table';
import chalk from 'chalk';

import { credentials } from 'utils/config';

class ConfigGet extends Command {
	async run() {
		const { flags } = this.parse(ConfigGet);

		try {
			const {
				name,
				email,
				apiKey,
				apiSecret,
				environment,
				telemetry,
				timeout,
			} = await credentials(this);

			if (flags.json) {
				this.log(JSON.stringify(await credentials(this)));
				this.exit(0);
			}

			const table = new Table();

			table.push(
				{
					[`${chalk.green.bold('Name')}`]: name,
				},
				{
					[`${chalk.green.bold('Email')}`]: email,
				},
				{
					[`${chalk.green.bold('API Key')}`]: apiKey,
				},
				{
					[`${chalk.green.bold('API Secret')}`]: apiSecret,
				},
				{
					[`${chalk.green.bold('Environment')}`]: environment,
				},
				{
					[`${chalk.green.bold('Telemetry')}`]: telemetry,
				},
				{
					[`${chalk.green.bold('Timeout(ms)')}`]: timeout,
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

ConfigGet.flags = {
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

ConfigGet.description = 'Outputs your user configuration.';

module.exports.ConfigGet = ConfigGet;
