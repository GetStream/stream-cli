import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid';

import { auth } from '../../utils/auth';
import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';
import { credentials } from '../../utils/config';

export class UserBan extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold('Name of the channel.'),
            default: uuid(),
            required: false,
        }),
        type: flags.string({
            char: 't',
            description: chalk.green.bold('Channel type.'),
            required: true,
        }),
        user: flags.string({
            char: 'u',
            description: chalk.green.bold('ID of the user.'),
            required: true,
        }),
        reason: flags.string({
            char: 'r',
            description: chalk.green.bold('Reason to place ban.'),
            required: false,
        }),
        timeout: flags.string({
            char: 't',
            description: chalk.green.bold('Duration in minutes.'),
            default: '60',
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(UserBan);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

            const payload = {};
            if (flags.timeout) payload.timeout = flags.timeout;
            if (flags.reason) payload.reason = flags.reason;

            await client.banUser(flags.user, payload);

            exit(
                `${flags.user} has been banned from ${flags.type}:${
                    flags.channel
                }`,
                {
                    emoji: 'warning',
                }
            );
        } catch (err) {
            apiError(err);
        }
    }
}

UserBan.description = 'Ban users from indefinitely or by a per-minute period.';
