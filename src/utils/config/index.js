import { StreamChat } from 'stream-chat';
import emoji from 'node-emoji';
import chalk from 'chalk';
import fs from 'fs-extra';
import path from 'path';

export async function credentials(config) {
    try {
        if (!(await fs.pathExists(config))) {
            await fs.outputJson(config, {
                apiKey: '',
                apiSecret: '',
            });
        }

        const { apiKey, apiSecret } = await fs.readJson(config);

        if (!apiKey.length || !apiSecret.length) {
            console.log(
                chalk.red(
                    `Credentials not found. Run ${chalk.bold(
                        'chat config:set'
                    )} to generate a configuration file. ${emoji.get(
                        'pensive'
                    )}`
                )
            );

            process.exit(0);
        }

        return { apiKey, apiSecret };
    } catch (err) {
        console.log(err);
        process.exit(1);
    }
}
