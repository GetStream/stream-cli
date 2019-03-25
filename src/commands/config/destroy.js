const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');
const fs = require('fs-extra');

class ConfigDestroy extends Command {
	async run() {
		const { flags } = this.parse(ConfigDestroy);

		try {
			await fs.remove(path.join(this.config.configDir, 'config.json'));

			if (!flags.force) {
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
			}

			this.log(
				`Config destroyed. Run the command ${chalk.bold(
					'stream config:set'
				)} to generate a new stream configuration file.`
			);

			this.exit(0);
		} catch (error) {
			this.error(error, { exit: 1 });
		}
	}
}

ConfigDestroy.flags = {
	force: flags.boolean({
		char: 'f',
		description: 'Force remove Stream configuration from cache.',
		required: false,
	}),
};

module.exports.ConfigDestroy = ConfigDestroy;
