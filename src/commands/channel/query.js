import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid/v4';

import { exit } from '../../utils/response';
import { authError, apiError } from '../../utils/error';
import { credentials } from '../../utils/config';

export class ChannelQuery extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold('ID of channel.'),
            default: uuid(),
            required: false,
        }),
        type: flags.string({
            char: 't',
            description: chalk.blue.bold('Type of channel.'),
            options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
            required: false,
        }),
        filter: flags.string({
            char: 'f',
            description: chalk.blue.bold('Filters to apply.'),
            required: false,
        }),
        sort: flags.string({
            char: 's',
            description: chalk.blue.bold('Sort to apply.'),
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(ChannelQuery);
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const { apiKey, apiSecret } = await credentials(config);
            if (!apiKey || !apiSecret) return authError();

            const client = new StreamChat(apiKey, apiSecret);

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
