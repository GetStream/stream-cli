const { Command, flags } = require('@oclif/command');
const Table = require('cli-table');
const chalk = require('chalk');
const jwt = require('jsonwebtoken');
const path = require('path');

const { credentials } = require('../../utils/config');

class DebugToken extends Command {
    async run() {
        const { apiKey, apiSecret } = await credentials(this);
        const { flags } = this.parse(DebugToken);

        try {
            const decoded = await jwt.verify(flags.jwt, apiSecret, {
                complete: true,
            });

            if (!decoded) {
                this.warn('Invalid JWT token or Stream API secret.');
                this.exit(0);
            }

            if (flags.raw) {
                this.log(decoded);
                this.exit();
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
    jwt: flags.string({
        char: 'j',
        description: chalk.blue.bold('The JWT token you are trying to debug.'),
        required: true,
    }),
    raw: flags.string({
        char: 'r',
        description: chalk.blue.bold(
            'A raw object containing the header, signature, and payload of your JWT.'
        ),
        required: false,
    }),
};

module.exports.DebugToken = DebugToken;
