const { Command, flags } = require('@oclif/command');
const Table = require('cli-table');
const chalk = require('chalk');

const { credentials } = require('../../utils/config');

class ConfigGet extends Command {
	async run() {
		const { flags } = this.parse(ConfigGet);

		try {
			const {
				name,
				email,
				apiKey,
				apiSecret,
				apiBaseUrl,
				environment,
				telemetry,
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
					[`${chalk.green.bold('API Base URL')}`]: apiBaseUrl,
				},
				{
					[`${chalk.green.bold('Environment')}`]: environment,
				},
				{
					[`${chalk.green.bold('Telemetry')}`]: telemetry,
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
