const { Command } = require('@oclif/command');
const { prompt } = require('enquirer');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');
const fs = require('fs-extra');

export class ConfigDestroy extends Command {
    async run() {
        try {
            await fs.remove(path.join(this.config.configDir, 'config.json'));

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

            this.log(
                `Config destroyed. Run the command ${chalk.bold(
                    'stream config:set'
                )} to generate a new config.`,
                emoji.get('rocket')
            );

            this.exit(0);
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ConfigDestroy.description = 'Destroy Stream configuration entirely.';
