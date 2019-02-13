import { Command, flags } from '@oclif/command';
import emoji from 'node-emoji';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid';

import { auth } from '../../utils/auth';

export class MessageRemove extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.green.bold('Channel ID.'),
            default: uuid(),
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(MessageRemove);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json'),
                this
            );

            await client.deleteMessage(flags.id);

            this.log(
                `The message ${flags.id} has been removed!`,
                emoji.get('wastebasket')
            );
            this.exit(0);
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

MessageRemove.description = 'Send messages to a channel.';
