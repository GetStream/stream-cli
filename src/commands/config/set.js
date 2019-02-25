const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const emoji = require('node-emoji');
const chalk = require('chalk');
const opn = require('opn');
const path = require('path');
const fs = require('fs-extra');

class ConfigSet extends Command {
    async run() {
        const { flags } = this.parse(ConfigSet);
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const exists = await fs.pathExists(config);

            let data = {};

            if (!flags.key && !flags.secret) {
                if (exists) {
                    const answer = await prompt({
                        type: 'confirm',
                        name: 'continue',
                        message: chalk.red.bold(
                            `This command will delete your current configuration. Are you sure you want to continue? ${emoji.get(
                                'warning'
                            )} `
                        ),
                    });

                    if (!answer.continue) {
                        this.exit(0);
                    }
                } else {
                    const answer = await prompt({
                        type: 'confirm',
                        name: 'continue',
                        message: `Do you have an existing account with Stream? If not, please enter "N".`,
                    });

                    if (!answer.continue) {
                        opn('https://getstream.io');

                        this.log(
                            chalk.yellow(
                                `Redirecting you to https://getstream.io`
                            ),
                            emoji.get('earth_americas')
                        );
                        this.exit(0);
                    }
                }

                data = await prompt([
                    {
                        type: 'input',
                        name: 'apiKey',
                        message: `What's your API key? ${emoji.get('lock')}`,
                    },
                    {
                        type: 'input',
                        name: 'apiSecret',
                        message: `What's your API secret? ${emoji.get('lock')}`,
                    },
                ]);
            } else {
                data.apiKey = flags.key;
                data.apiSecrety = flags.secret;
            }

            await fs.writeJson(config, {
                apiKey: flags.key ? flags.key : data.apiKey,
                apiSecret: flags.secret ? flags.secret : data.apiSecret,
            });

            this.log(
                chalk.bold(`Your config has been generated!`),
                emoji.get('rocket')
            );
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ConfigSet.flags = {
    key: flags.string({
        char: 'k',
        description: chalk.blue.bold('API key for config.'),
    }),
    secret: flags.string({
        char: 's',
        description: chalk.blue.bold('API secret for config.'),
    }),
};

module.exports.ConfigSet = ConfigSet;
