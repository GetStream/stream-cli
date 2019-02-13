import { Command, flags } from '@oclif/command';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid';

import { auth } from '../../utils/auth';
import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';

export class UserBan extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.green.bold('Channel ID.'),
            default: uuid(),
            required: true,
        }),
        type: flags.string({
            char: 't',
            description: chalk.green.bold('Type of channel.'),
            options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
            required: true,
        }),
        user: flags.string({
            char: 'u',
            description: chalk.green.bold('User ID.'),
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
                path.join(this.config.configDir, 'config.json'),
                this
            );

            const payload = {};
            if (flags.timeout) payload.timeout = flags.timeout;
            if (flags.reason) payload.reason = flags.reason;

            await client.banUser(flags.user, payload);

            exit(
                `${flags.user} has been banned from ${flags.type}:${flags.id}`,
                {
                    emoji: 'warning',
                },
                this
            );
        } catch (err) {
            apiError(err, this);
        }
    }
}

UserBan.description = 'Ban users indefinitely or by a per-minute period.';
