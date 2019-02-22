const Command = require('@oclif/command');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');
const fs = require('fs');

export class ConfigDestroy extends Command {
    async run() {
        try {
            await fs.remove(path.join(this.config.configDir, 'config.json'));

            this.log(
                `Config destroyed. Run the command ${chalk.blue.bold(
                    'config:set'
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
