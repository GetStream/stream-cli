import { Command, flags } from '@oclif/command';
import emoji from 'node-emoji';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid';

import { auth } from '../../utils/auth';

export class MessageSend extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold('Channel ID.'),
            default: uuid(),
            required: false,
        }),
        user: flags.string({
            char: 'u',
            description: chalk.blue.bold('ID of user.'),
            default: '*',
            required: true,
        }),
        type: flags.string({
            char: 't',
            description: chalk.blue.bold('Type of channel.'),
            options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
            required: true,
        }),
        message: flags.string({
            char: 'm',
            description: chalk.blue.bold('Message to send.'),
            required: true,
        }),
        attachments: flags.string({
            char: 'a',
            description: chalk.blue.bold('JSON payload of attachments'),
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(MessageSend);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json'),
                this
            );

            await client.updateUser({
                id: flags.uid,
                role: 'admin',
            });

            await client.setUser({ id: flags.uid, status: 'invisible' });
            const channel = client.channel(flags.type, flags.channel);

            const payload = {
                text: flags.message,
            };

            if (flags.attachments) {
                payload.attachments = JSON.parse(flags.attachments);
            }

            await channel.sendMessage(payload);

            const message = chalk.blue(
                `Message ${chalk.bold(
                    flags.message
                )} has been sent to the ${chalk.bold(
                    flags.channel
                )} channel by ${chalk.bold(flags.uid)}!`
            );

            this.log(message, emoji.get('smile'));
            this.exit();
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

MessageSend.description = 'Send messages to a channel';
