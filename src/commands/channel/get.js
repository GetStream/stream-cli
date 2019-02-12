import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
import stringify from 'json-stringify-pretty-compact';
import cardinal from 'cardinal';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';

import { auth } from '../../utils/auth';
import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';

export class ChannelGet extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold('Channel ID.'),
            required: false,
        }),
        type: flags.string({
            char: 't',
            description: chalk.blue.bold('Type of channel.'),
            options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
            required: false,
        }),
        config: flags.boolean({
            char: 'c',
            description: chalk.blue.bold('Return channel config values only.'),
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(ChannelGet);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );
            const channel = await client.queryChannels(
                { id: flags.id, type: flags.type },
                { last_message_at: -1 },
                {
                    subscribe: false,
                }
            );

            const payload = cardinal.highlight(
                stringify(
                    flags.config ? channel[0].data.config : channel[0].data,
                    {
                        maxLength: 100,
                    },
                    {
                        linenos: true,
                    }
                )
            );

            exit(payload, { newline: true });
        } catch (err) {
            apiError(err);
        }
    }
}

ChannelGet.description = 'Get a channel.';
