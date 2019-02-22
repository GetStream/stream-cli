const { Command } = require('@oclif/command');
const prompt = require('enquirer');
const emoji = require('node-emoji');
const chalk = require('chalk');
const opn = require('opn');
const path = require('path');
const fs = require('fs-extra');

export class ConfigSet extends Command {
    async run() {
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const exists = await fs.pathExists(config);

            if (exists) {
                const answer = await prompt({
                    type: 'confirm',
                    name: 'continue',
                    message: chalk.yellow.bold(
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
                    message: chalk.green(
                        `Do you have an existing account with Stream? If not, please enter "N". ${emoji.get(
                            'rocket'
                        )} `
                    ),
                });

                if (!answer.continue) {
                    opn('https://getstream.io');

                    this.log(
                        chalk.yellow(`Redirecting you to https://getstream.io`),
                        emoji.get('earth_americas')
                    );
                    this.exit(0);
                }
            }

            const data = await prompt([
                {
                    type: 'input',
                    name: 'apiKey',
                    message: chalk.green(
                        `What's your API key? ${emoji.get('lock')}`
                    ),
                },
                {
                    type: 'input',
                    name: 'apiSecret',
                    message: chalk.green(
                        `What's your API secret? ${emoji.get('lock')}`
                    ),
                },
            ]);

            await fs.writeJson(config, {
                apiKey: data.apiKey,
                apiSecret: data.apiSecret,
            });

            this.log(
                chalk.green(`Your config has been generated!`),
                emoji.get('rocket')
            );
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ConfigSet.description = 'Initialize the config with your Stream credentials.';
