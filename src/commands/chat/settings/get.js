const { Command, flags } = require('@oclif/command');
const Table = require('cli-table');
const chalk = require('chalk');

const { credentials } = require('../../../utils/config');

class SettingsGet extends Command {
    async run() {
        const { flags } = this.parse(ReactionRemove);

        try {
            const { name, email, apiKey, apiSecret } = await credentials(this);

            this.exit(0);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

SettingsGet.flags = {
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.SettingsGet = SettingsGet;
