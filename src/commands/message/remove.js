import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import fs from 'fs-extra';
import uuid from 'uuid';

import { auth } from '../../utils/auth';
import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';

export class MessageRemove extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold('Channel ID.'),
            default: uuid(),
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(MessageRemove);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

            await client.deleteMessage(flags.id);

            exit(`The message ${flags.id} has been removed!`, {
                emoji: 'wastebasket',
            });
        } catch (err) {
            apiError(err);
        }
    }
}

MessageRemove.description = 'Send messages to a channel.';
