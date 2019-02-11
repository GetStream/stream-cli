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

export class ChannelGet extends Command {
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
    };

    async run() {
        const { flags } = this.parse(ChannelGet);
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const { apiKey, apiSecret } = await credentials(config);
            if (!apiKey || !apiSecret) return authError();

            const client = new StreamChat(apiKey, apiSecret);

            const filter = { members: { $in: ['Nick'] } };
            const sort = { last_message_at: -1 };

            const channels = await client.queryChannels(filter, sort, {
                watch: true,
                state: true,
            });

            channels.map(channel => {
                console.log(channel);
            });

            // const message = chalk.blue(
            //     `The channel ${flags.name} has been modified!`
            // );
            //
            // console.log(`${timestamp}:`, message, emoji.get('rocket'));

            this.exit(0);
        } catch (err) {
            apiError(err);
        }
    }
}

ChannelGet.description = 'Edit a channel.';
