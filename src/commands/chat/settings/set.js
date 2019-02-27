const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const axios = require('axios');

const { credentials } = require('../../../utils/config');

class SettingsSet extends Command {
    async run() {
        const { flags } = this.parse(SettingsSet);

        try {
            const { apiKey, apiSecret } = await credentials(this);

            if (!flags.name || !flags.p12) {
                const res = await prompt([
                    {
                        type: 'input',
                        name: 'name',
                        message: `What is your name?`,
                        required: true,
                    },
                ]);

                for (const key in res) {
                    if (res.hasOwnProperty(key)) {
                        flags[key] = res[key];
                    }
                }

                const setting = null;

                if (flags.json) {
                    this.log(settings);
                    this.exit(0);
                }
            }

            this.log('Your Stream settings have been updated!');
            this.exit(0);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

SettingsSet.flags = {
    name: flags.string({
        char: 'n',
        description: 'Full name for settings.',
        required: false,
    }),
    p12: flags.string({
        char: 'p12',
        description: 'A .p12 file for push notifications.',
        required: false,
    }),
    pem: flags.string({
        char: 'pem',
        description: 'The .pem file associated with your .p12',
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.SettingsSet = SettingsSet;
