import { Command, flags } from '@oclif/command';
import emoji from 'node-emoji';
import chalk from 'chalk';
import path from 'path';
import fs from 'fs-extra';

import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';
import { credentials } from '../../utils/config';

export class ConfigDestroy extends Command {
    async run() {
        const config = path.join(this.config.configDir, 'config.json');

        try {
            await fs.remove(config);
            exit(`Config destroyed...`, { emoji: 'cry' });
        } catch (err) {
            authError(err);
        }
    }
}

ConfigDestroy.description = 'Destroys CLI config settings.';
