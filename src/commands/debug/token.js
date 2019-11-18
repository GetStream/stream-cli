const { Command, flags } = require('@oclif/command');
const Table = require('cli-table');
const { prompt } = require('enquirer');
const chalk = require('chalk');
const jwt = require('jsonwebtoken');

const { credentials } = require('../../utils/config');

class DebugToken extends Command {
	async run() {
		const { flags } = this.parse(DebugToken);

		try {
			if (!flags.token) {
				const res = await prompt([
					{
						type: 'input',
						name: 'jwt',
						message: `What is the Stream token you would like to debug?`,
						required: true,
					},
				]);

				flags.jwt = res.jwt;
			}

			const { apiSecret } = await credentials(this);

			const decoded = await jwt.verify(flags.jwt, apiSecret, {
				complete: true,
			});

			if (flags.json) {
				this.log(JSON.stringify(decoded));
				this.exit(0);
			}

			const table = new Table();

			table.push(
				{
					[`${chalk.green.bold('Header Type')}`]: decoded.header.typ,
				},
				{
					[`${chalk.green.bold('Header Algorithm')}`]: decoded.header
						.alg,
				},
				{
					[`${chalk.green.bold('Signature')}`]: decoded.signature,
				},
				{
					[`${chalk.green.bold('User ID')}`]: decoded.payload.user_id,
				}
			);

			this.log(table.toString());
			this.exit(0);
		} catch (error) {
			this.error('Malformed JWT token or Stream API secret.', {
				exit: 1,
			});
		}
	}
}

DebugToken.flags = {
	token: flags.string({
		char: 't',
		description: 'The Stream token you are trying to debug.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

DebugToken.description = 'Debugs a JWT token provided by Stream.';

module.exports.DebugToken = DebugToken;
