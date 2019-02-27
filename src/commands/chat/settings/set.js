const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const axios = require('axios');

const { credentials } = require('../../../utils/config');

class SettingsSet extends Command {
    async run() {
        const { apiKey, apiSecret } = await credentials(this);
        const { flags } = this.parse(SettingsSet);

        try {
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
        char: 'p',
        description: 'A .p12 file for push notifications.',
        required: false,
    }),
};

module.exports.SettingsSet = SettingsSet;
