import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import fs from 'fs-extra';

import { exit } from '../../utils/response';
import { authError } from '../../utils/error';
import { credentials } from '../../utils/config';

export class ModerateBan extends Command {
    static flags = {
        user: flags.string({
            char: 'u',
            description: chalk.blue.bold('ID of user.'),
            exclusive: ['message'],
            required: true,
        }),
        reason: flags.string({
            char: 'r',
            description: chalk.blue.bold('Reason for timeout.'),
            required: true,
        }),
        timeout: flags.string({
            char: 't',
            description: chalk.blue.bold('Timeout in minutes.'),
            default: '60',
            required: true,
        }),
    };

    async run() {
        const { flags } = this.parse(ModerateBan);
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const { apiKey, apiSecret } = await credentials(config);
            if (!apiKey || !apiSecret) return authError();

            const client = new StreamChat(apiKey, apiSecret);

            const timestamp = chalk.yellow.bold(
                moment().format('dddd, MMMM Do YYYY [at] h:mm:ss A')
            );

            await client.banUser(flags.user, {
                timeout: Number(flags.timeout),
                reason: flags.reason,
            });

            exit(`The user ${flags.user} has been banned!`, 'banned');
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ModerateBan.description = 'Flag users and messages.';
