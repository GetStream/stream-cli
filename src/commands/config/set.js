import { Command, flags } from '@oclif/command';
import { prompt } from 'enquirer';
import emoji from 'node-emoji';
import chalk from 'chalk';
import path from 'path';
import fs from 'fs-extra';

import { exit } from '../../utils/response';
import { authError } from '../../utils/error';
import { credentials } from '../../utils/config';

export class ConfigSet extends Command {
    async run() {
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const { apiKey, apiSecret } = await credentials(config);

            if (fs.existsSync(config)) {
                const answer = await prompt({
                    type: 'confirm',
                    name: 'continue',
                    message: chalk.yellow.bold(
                        `This command will delete your current configuration. Do you want to continue? ${emoji.get(
                            'warning'
                        )} `
                    ),
                });

                if (!answer.continue) {
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

            exit(`Your config has been saved!`, 'rocket');
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ConfigSet.description = 'Stores your Stream API key and secret.';
