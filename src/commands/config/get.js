const { Command } = require('@oclif/command');
const Table = require('cli-table');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');

const { credentials } = require('../../utils/config');

class ConfigGet extends Command {
    async run() {
        const { name, email, apiKey, apiSecret } = await credentials(this);

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
            }
        );

        this.log(table.toString());
        this.exit(0);
    }
}

module.exports.ConfigGet = ConfigGet;
