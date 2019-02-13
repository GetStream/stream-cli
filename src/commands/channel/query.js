import { Command, flags } from '@oclif/command';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid/v4';

import { auth } from '../../utils/auth';

export class ChannelQuery extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.green.bold('Channel ID.'),
            default: uuid(),
            required: false,
        }),
        type: flags.string({
            char: 't',
            description: chalk.green.bold('Type of channel.'),
            options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
            required: false,
        }),
        filter: flags.string({
            char: 'f',
            description: chalk.green.bold('Filters to apply.'),
            required: false,
        }),
        sort: flags.string({
            char: 's',
            description: chalk.green.bold('Sort to apply.'),
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(ChannelQuery);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json'),
                this
            );

            const filter = flags.filters ? JSON.parse(flags.filters) : {};
            const sort = flags.sort ? JSON.parse(flags.sort) : {};

            const channels = await client.queryChannels(filter, sort, {
                subscribe: false,
            });

            this.log(channels[0]);

            this.exit(0);
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ChannelQuery.description = 'Query a specific hannel.';
