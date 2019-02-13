import { Command, flags } from '@oclif/command';
import emoji from 'node-emoji';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid/v4';

import { auth } from '../../utils/auth';

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
                path.join(this.config.configDir, 'config.json'),
                this
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

            this.log(
                `The channel ${flags.id} has been modified!`,
                emoji.get('rocket')
            );
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ChannelEdit.description = 'Edit a channel.';
