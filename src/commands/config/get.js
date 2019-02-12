import { Command, flags } from '@oclif/command';
import stringify from 'json-stringify-pretty-compact';
import cardinal from 'cardinal';
import emoji from 'node-emoji';
import chalk from 'chalk';
import path from 'path';
import fs from 'fs-extra';

import { authError } from '../../utils/error';
import { credentials } from '../../utils/config';

export class ConfigGet extends Command {
    async run() {
        const config = path.join(this.config.configDir, 'config.json');
        const { apiKey, apiSecret } = await credentials(config);

        if (apiKey && apiSecret) {
            const creds = cardinal.highlight(stringify({ apiKey, apiSecret }), {
                linenos: true,
            });

            exit(creds, { newline: true });
        } else {
            return authError();
        }
    }
}

ConfigGet.description = 'Retrieves config credentials for CLI.';
