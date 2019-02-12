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

export class ModerateMute extends Command {
    static flags = {
        user: flags.string({
            char: 'u',
            description: chalk.green.bold('ID of user.'),
            required: true,
        }),
    };

    async run() {
        const { flags } = this.parse(ModerateMute);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

            const timestamp = chalk.yellow.bold(
                moment().format('dddd, MMMM Do YYYY [at] h:mm:ss A')
            );

            await authClient.muteUser(flags.user);

            exit(`The message ${flags.user} has been flagged!`, {
                emoji: 'two_flags',
            });
        } catch (err) {
            apiError(err);
        }
    }
}

ModerateMute.description = 'Flag users and messages.';
