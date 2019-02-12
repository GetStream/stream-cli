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

export class ChannelEdit extends Command {
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
            required: true,
        }),
        name: flags.string({
            char: 'n',
            description: chalk.blue.bold('Name of room.'),
            required: true,
        }),
        url: flags.string({
            char: 'u',
            description: chalk.blue.bold('URL to channel image.'),
            required: false,
        }),
        members: flags.string({
            char: 'm',
            description: chalk.blue.bold('Comma separated list of members.'),
            required: false,
        }),
        reason: flags.string({
            char: 'r',
            description: chalk.blue.bold('Reason for changing channel.'),
            required: false,
        }),
        data: flags.string({
            char: 'd',
            description: chalk.blue.bold('Additional data as a JSON payload.'),
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(ChannelEdit);
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const { apiKey, apiSecret } = await credentials(config);
            if (!apiKey || !apiSecret) return authError();

            const client = new StreamChat(apiKey, apiSecret);

            const timestamp = chalk.yellow.bold(
                moment().format('dddd, MMMM Do YYYY [at] h:mm:ss A')
            );

            const channel = await client.channel(flags.type, flags.channel);

            const payload = {};
            if (flags.url) payload.url = flags.url;
            if (flags.name) payload.name = flags.name;

            await channel.update(payload, {
                text: flags.reason,
            });

            exit(`The channel ${flags.name} has been modified!`, {
                emoji: 'rocket',
            });
        } catch (err) {
            apiError(err);
        }
    }
}

ChannelEdit.description = 'Edit a channel.';
