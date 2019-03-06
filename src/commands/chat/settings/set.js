const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');

class SettingsSet extends Command {
    async run() {
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

                if (flags.json) {
                    this.log('settings');
                    this.exit(0);
                }
            }

            this.log('Your Stream orginzation settings have been updated.');
            this.exit(0);
        } catch (error) {
            this.error(error || 'A Stream CLI error has occurred.', {
                exit: 1,
            });
        }
    }
}

SettingsSet.flags = {
    p12: flags.string({
        char: 'p',
        description: '.p12 file.',
        required: false,
    }),
    key: flags.string({
        char: 'k',
        description: '.p8 file',
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
