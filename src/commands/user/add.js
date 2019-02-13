import { Command, flags } from '@oclif/command';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid';

import { auth } from '../../utils/auth';
import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';

export class UserAdd extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.green.bold('Channel name.'),
            default: uuid(),
            required: true,
        }),
        type: flags.string({
            char: 't',
            description: chalk.green.bold('Channel type.'),
            required: true,
        }),
        moderators: flags.string({
            char: 'm',
            description: chalk.green.bold(
                'Comma separated list of moderators to add.'
            ),
            required: true,
        }),
    };

    async run() {
        const { flags } = this.parse(UserAdd);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json'),
                this
            );

            const channel = await client.channel(flags.type, flags.id);
            await channel.addModerators(flags.moderators.split(','));

            exit(
                `${flags.moderators} have been added as moderators to channel ${
                    flags.type
                }:${flags.id}`,
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

UserAdd.description = 'Remove users from a channel.';
