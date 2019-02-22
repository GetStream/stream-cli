const { Command } = require('@oclif/command');
const Table = require('cli-table');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');

const { credentials } = require('../../utils/config');

class ConfigGet extends Command {
    async run() {
        const config = path.join(this.config.configDir, 'config.json');
        const { apiKey, apiSecret } = await credentials(config, this);

        if (apiKey && apiSecret) {
            const table = new Table({
                head: [
                    chalk.green.bold('API Key'),
                    chalk.green.bold('API Secret'),
                ],
                colWidths: [25, 70],
            });

            table.push([apiKey, apiSecret]);

            this.log(table.toString());
            this.exit(0);
        } else {
            this.error(
                `Credentials not found. Run ${chalk.bold(
                    'stream config:set'
                )} to generate a configuration file. ${emoji.get('caution')}`,
                { exit: 1 }
            );
        }
    }
}

module.exports.ConfigGet = ConfigGet;
