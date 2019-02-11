import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';

import { exit } from '../../utils/response';
import { authError } from '../../utils/error';
import { credentials } from '../../utils/config';

export class ChannelGet extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold('Name of the channel.'),
            default: uuid(),
            required: false,
        }),
        type: flags.string({
            char: 't',
            description: chalk.blue.bold('Type of the channel.'),
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

            const timestamp = chalk.yellow.bold(
                moment().format('dddd, MMMM Do YYYY [at] h:mm:ss A')
            );

            const filter = { members: { $in: ['Nick Tarsons'] } };
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
            this.error(err, { exit: 1 });
        }
    }
}

ChannelGet.description = 'Edit a channel.';
