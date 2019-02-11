import { Command, flags } from '@oclif/command';
import emoji from 'node-emoji';
import chalk from 'chalk';
import path from 'path';
import fs from 'fs-extra';

import { exit } from '../../utils/response';
import { credentials } from '../../utils/config';

export class ConfigDestroy extends Command {
    async run() {
        const config = path.join(this.config.configDir, 'config.json');

        try {
            await fs.remove(config);

            exit(`Config destroyed...`, 'cry');
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ConfigDestroy.description = 'Destroys CLI config settings.';
