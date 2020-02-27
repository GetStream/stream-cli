import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import emoji from 'node-emoji';
import chalk from 'chalk';
import path from 'path';
import fs from 'fs-extra';

class ConfigDestroy extends Command {
	async run() {
		const { flags } = this.parse(ConfigDestroy);

		const config = path.join(this.config.configDir, 'config.json');

		try {
			if (!flags.force) {
				const answer = await prompt({
					type: 'confirm',
					name: 'continue',
					message: chalk.red.bold(
						`This command will delete your current configuration. Are you sure you want to continue? ${emoji.get(
							'warning'
						)} `
					)
				});

				if (answer.continue) {
					await fs.remove(config);
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
		required: false
	})
};

ConfigDestroy.description = 'Destroys your user configuration.';

module.exports.ConfigDestroy = ConfigDestroy;
