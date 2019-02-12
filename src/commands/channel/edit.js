import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
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
            description: chalk.blue.bold('Channel ID.'),
            required: true,
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
        reason: flags.string({
            char: 'r',
            description: chalk.blue.bold('Reason for changing channel.'),
            required: true,
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
        const { flags } = this.parse(ChannelEdit);
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const { apiKey, apiSecret } = await credentials(config);
            if (!apiKey || !apiSecret) return authError();

            const client = new StreamChat(apiKey, apiSecret);

            const channel = await client.channel(flags.type, flags.id);

            let payload = {
                name: flags.name,
                updated_by: {
                    id: uuid(),
                    name: 'CLI',
                },
            };
            if (flags.image) payload.image = flags.image;
            if (flags.members) payload.members = flags.members.split(',');

            if (flags.data) {
                const parsed = JSON.parse(flags.data);
                payload = Object.assign({}, payload, parsed);
            }

            await channel.update(payload, {
                name: flags.name,
                text: flags.reason,
            });

            exit(`The channel ${flags.id} has been modified!`, {
                emoji: 'rocket',
            });
        } catch (err) {
            apiError(err);
        }
    }
}

ChannelEdit.description = 'Edit a channel.';
