import { Command, flags } from '@oclif/command';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid';

import { auth } from '../../utils/auth';
import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';
import { credentials } from '../../utils/config';

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
                'Comma separated list of moderators to remove.'
            ),
            required: true,
        }),
    };

    async run() {
        const { flags } = this.parse(UserAdd);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

            const channel = await client.channel(flags.type, flags.id);
            await channel.demoteModerators(flags.moderators.split(','));

            exit(
                `${
                    flags.moderators
                } have been removed as moderators from channel ${flags.type}:${
                    flags.id
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

UserAdd.description = 'Remove users from a channel.';
