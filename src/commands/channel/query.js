import { Command, flags } from '@oclif/command';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid/v4';

import { auth } from '../../utils/auth';
import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';

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
                path.join(this.config.configDir, 'config.json')
            );

            const filter = flags.filters ? JSON.parse(flags.filters) : {};
            const sort = flags.sort ? JSON.parse(flags.sort) : {};

            const channels = await client.queryChannels(filter, sort, {
                subscribe: false,
            });

            console.log(channels[0]);

            process.exit(0);
        } catch (err) {
            apiError(err);
        }
    }
}

ChannelQuery.description = 'Query a channel.';
