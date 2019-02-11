import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid/v4';

import { exit } from '../../utils/response';
import { authError } from '../../utils/error';
import { credentials } from '../../utils/config';

export class ChannelInit extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold('Name of channel.'),
            default: uuid(),
            required: false,
        }),
        type: flags.string({
            char: 't',
            description: chalk.blue.bold('Type of channel.'),
            options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
            required: true,
        }),
        name: flags.string({
            char: 'n',
            description: chalk.blue.bold('Name of room.'),
            required: true,
        }),
        image: flags.string({
            char: 'u',
            description: chalk.blue.bold('URL to channel image.'),
            required: false,
        }),
        members: flags.string({
            char: 'm',
            description: chalk.blue.bold('Comma separated list of members.'),
            required: false,
        }),
        data: flags.string({
            char: 'd',
            description: chalk.blue.bold('Additional data as a JSON payload.'),
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(ChannelInit);
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const { apiKey, apiSecret } = await credentials(config);
            if (!apiKey || !apiSecret) return authError();

            const client = new StreamChat(apiKey, apiSecret);

            const timestamp = chalk.yellow.bold(
                moment().format('dddd, MMMM Do YYYY [at] h:mm:ss A')
            );

            const payload = {
                name: flags.name,
            };
            if (flags.image) payload.image = flags.image;
            if (flags.members) payload.members = flags.members.split(',');
            if (flags.metadata) payload.data = JSON.parse(flags.metadata);

            const channel = await client.channel(
                flags.type,
                flags.channel,
                payload
            );

            exit(`The channel ${flags.name} has been initialized!`, 'rocket');
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ChannelInit.description = 'Initialize a channel.';
