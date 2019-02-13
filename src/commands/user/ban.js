import { Command, flags } from '@oclif/command';
import emoji from 'node-emoji';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid';

import { auth } from '../../utils/auth';

export class UserBan extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold('Channel ID.'),
            default: uuid(),
            required: true,
        }),
        type: flags.string({
            char: 't',
            description: chalk.blue.bold('Type of channel.'),
            options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
            required: true,
        }),
        user: flags.string({
            char: 'u',
            description: chalk.blue.bold('User ID.'),
            required: true,
        }),
        reason: flags.string({
            char: 'r',
            description: chalk.blue.bold('Reason to place ban.'),
            required: false,
        }),
        timeout: flags.string({
            char: 't',
            description: chalk.blue.bold('Duration in minutes.'),
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

            this.log(
                `${flags.user} has been added banned from channel ${
                    flags.type
                }:${flags.id}`,
                emoji.get('banned')
            );
            this.exit(0);
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

UserBan.description = 'Ban users indefinitely or by a per-minute period.';
