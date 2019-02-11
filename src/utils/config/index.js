import { StreamChat } from 'stream-chat';
import chalk from 'chalk';
import fs from 'fs-extra';
import path from 'path';

export async function credentials(config) {
    try {
        const { apiKey, apiSecret } = await fs.readJson(config);

        if (!apiKey || !apiSecret) {
            console.log(
                chalk.red(
                    `Credentials not found. Run ${chalk.bold(
                        'chat init'
                    )} to generate a configuration file. ${emoji.get(
                        'pensive'
                    )}`
                )
            );

            process.exit(0);
        }

        return { apiKey, apiSecret };
    } catch (err) {
        console.error(err);
        process.exit(1);
    }
}
