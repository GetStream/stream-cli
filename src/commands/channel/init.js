import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid/v4';

import { auth } from '../../utils/auth';
import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';

export class ChannelInit extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold('Channel ID.'),
            default: uuid(),
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

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

            let payload = {
                name: flags.name,
                created_by: {
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

            const channel = await client.channel(flags.type, flags.id, payload);
            await channel.create();

            exit(`The channel ${flags.name} has been initialized!`, {
                emoji: 'rocket',
            });
        } catch (err) {
            apiError(err);
        }
    }
}

ChannelInit.description = 'Initialize a channel.';
