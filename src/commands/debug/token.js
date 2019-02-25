const { Command, flags } = require('@oclif/command');
const Table = require('cli-table');
const chalk = require('chalk');
const jwt = require('jsonwebtoken');
const path = require('path');

const { credentials } = require('../../utils/config');

class DebugToken extends Command {
    async run() {
        const config = path.join(this.config.configDir, 'config.json');
        const { apiKey, apiSecret } = await credentials(config, this);
        const { flags } = this.parse(DebugToken);

        try {
            const decoded = await jwt.verify(flags.token, apiSecret, {
                complete: true,
            });

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
            this.error(err, { exit: 1 });
        }
    }
}

DebugToken.flags = {
    token: flags.string({
        char: 't',
        description: chalk.blue.bold('The JWT token you are trying to debug.'),
        required: true,
    }),
};

module.exports.DebugToken = DebugToken;
