import { Command, flags } from '@oclif/command';
import emoji from 'node-emoji';
import chalk from 'chalk';
import path from 'path';

import { auth } from '../../utils/auth';

export class UserRemove extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold('Channel name.'),
            required: true,
        }),
        type: flags.string({
            char: 't',
            description: chalk.blue.bold('Channel type.'),
            required: true,
        }),
        moderators: flags.string({
            char: 'm',
            description: chalk.blue.bold(
                'Comma separated list of moderators to remove.'
            ),
            required: true,
        }),
    };

    async run() {
        const { flags } = this.parse(UserRemove);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json'),
                this
            );

            const channel = await client.channel(flags.type, flags.id);
            await channel.demoteModerators(flags.moderators.split(','));

            this.log(
                `${flags.moderators} have been removed as moderators from the ${
                    flags.type
                } channel ${flags.id}`,
                emoji.get('warning')
            );
            this.exit(0);
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

UserRemove.description = 'Remove users from a channel.';
