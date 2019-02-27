const { Command, flags } = require('@oclif/command');
const Table = require('cli-table');
const chalk = require('chalk');
const jwt = require('jsonwebtoken');
const path = require('path');

const { credentials } = require('../../utils/config');

class DebugToken extends Command {
    async run() {
        const { flags } = this.parse(DebugToken);

        try {
            const { apiKey, apiSecret } = await credentials(this);

            const decoded = await jwt.verify(flags.jwt, apiSecret, {
                complete: true,
            });

            if (!decoded) {
                this.warn('Invalid JWT token or Stream API secret.');
                this.exit(0);
            }

            if (flags.json) {
                this.log(decoded);
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
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

DebugToken.flags = {
    token: flags.string({
        char: 't',
        description: 'The Stream token you are trying to debug.',
        required: true,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.DebugToken = DebugToken;
