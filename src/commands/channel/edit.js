import { Command, flags } from '@oclif/command';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid/v4';

import { auth } from '../../utils/auth';
import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';

export class ChannelEdit extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.green.bold('Channel ID.'),
            required: true,
        }),
        type: flags.string({
            char: 't',
            description: chalk.green.bold('Type of channel.'),
            options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
            required: true,
        }),
        name: flags.string({
            char: 'n',
            description: chalk.green.bold('Name of room.'),
            required: true,
        }),
        url: flags.string({
            char: 'u',
            description: chalk.green.bold('URL to channel image.'),
            required: false,
        }),
        reason: flags.string({
            char: 'r',
            description: chalk.green.bold('Reason for changing channel.'),
            required: true,
        }),
        members: flags.string({
            char: 'm',
            description: chalk.green.bold('Comma separated list of members.'),
            required: false,
        }),
        data: flags.string({
            char: 'd',
            description: chalk.green.bold('Additional data as a JSON payload.'),
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(ChannelEdit);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );
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
