import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import fs from 'fs-extra';

import { exit } from '../../utils/response';
import { authError } from '../../utils/error';
import { credentials } from '../../utils/config';

export class ModerateMute extends Command {
    static flags = {
        user: flags.string({
            char: 'u',
            description: chalk.blue.bold('ID of user.'),
            required: true,
        }),
    };

    async run() {
        const { flags } = this.parse(ModerateMute);
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const { apiKey, apiSecret } = await credentials(config);
            if (!apiKey || !apiSecret) return authError();

            const client = new StreamChat(apiKey, apiSecret);

            const timestamp = chalk.yellow.bold(
                moment().format('dddd, MMMM Do YYYY [at] h:mm:ss A')
            );

            await authClient.muteUser(flags.user);

            exit(`The message ${flags.user} has been flagged!`, 'two_flags');
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ModerateMute.description = 'Flag users and messages.';
