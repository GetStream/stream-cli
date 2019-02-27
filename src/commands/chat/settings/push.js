const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const axios = require('axios');

const { credentials } = require('../../../utils/config');

class SettingsPush extends Command {
    async run() {
        const { flags } = this.parse(SettingsPush);

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

            this.log('Your push notification settings have been updated.');
            this.exit(0);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

SettingsPush.flags = {
    type: flags.boolean({
        char: 't',
        description: 'Type of configuration.',
        options: ['apn', 'firebase', 'webhook'],
        required: false,
    }),
    auth_key: flags.string({
        char: 'a',
        description: 'Private auth key for APN.',
        required: false,
    }),
    key_id: flags.string({
        char: 'k',
        description: 'Key ID for APN.',
        required: false,
    }),
    team_id: flags.string({
        char: 't',
        description: 'Team ID for APN.',
        required: false,
    }),
    pem_cert: flags.string({
        char: 'p',
        description: 'Private RSA key for APN (.pem).',
        required: false,
    }),
    p12_cert: flags.string({
        char: 'b',
        description: 'Base64 encoded .p12 file for APN.',
        required: false,
    }),
    notification_template: flags.string({
        char: 'n',
        description: 'Interpolated JSON template for notifications.',
        required: false,
    }),
    api_key: flags.string({
        char: 'a',
        description: 'API key for Firebase.',
        required: false,
    }),
    webhook_url: flags.string({
        char: 'w',
        description: 'Webhook URL to receive notifications to.',
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.SettingsPush = SettingsPush;
