import { Command, flags } from '@oclif/command';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import fs from 'fs-extra';

import { auth } from '../../utils/auth';
import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';
import { credentials } from '../../utils/config';

export class ModerateFlag extends Command {
    static flags = {
        user: flags.string({
            char: 'u',
            description: chalk.green.bold('ID of user.'),
            exclusive: ['message'],
            required: false,
        }),
        message: flags.string({
            char: 'm',
            description: chalk.green.bold('ID of message.'),
            exclusive: ['user'],
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(ModerateFlag);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

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

            exit(message, { emoji: 'crossed_flags' });
        } catch (err) {
            apiError(err);
        }
    }
}

ModerateFlag.description = 'Flag users and messages.';
