import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import fs from 'fs-extra';

import { exit } from '../../utils/response';
import { authError, apiError } from '../../utils/error';
import { credentials } from '../../utils/config';

export class ModerateFlag extends Command {
    static flags = {
        user: flags.string({
            char: 'u',
            description: chalk.blue.bold('ID of user.'),
            exclusive: ['message'],
            required: false,
        }),
        message: flags.string({
            char: 'm',
            description: chalk.blue.bold('ID of message.'),
            exclusive: ['user'],
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(ModerateFlag);
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const { apiKey, apiSecret } = await credentials(config);
            if (!apiKey || !apiSecret) return authError();

            const client = new StreamChat(apiKey, apiSecret);

            const timestamp = chalk.yellow.bold(
                moment().format('dddd, MMMM Do YYYY [at] h:mm:ss A')
            );

            if (flags.user) {
                await client.flagUser(flags.user);

                const message = chalk.blue(
                    `The user ${flags.user} has been flagged!`
                );
            } else if (flags.message) {
                await client.flagMessage(flags.message);

                const message = chalk.blue(
                    `The message ${flags.user} has been flagged!`
                );
            } else {
                console.log(chalk.red.bold(`Please pass a valid command.`));

                this.exit(0);
            }

            exit(message, 'crossed_flags');
        } catch (err) {
            apiError(err);
        }
    }
}

ModerateFlag.description = 'Flag users and messages.';
