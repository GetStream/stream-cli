import { Command, flags } from '@oclif/command';
import stringify from 'json-stringify-pretty-compact';
import cardinal from 'cardinal';
import chalk from 'chalk';
import path from 'path';

import { auth } from '../../utils/auth';

export class ChannelGet extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.green.bold('Channel ID.'),
            required: false,
        }),
        type: flags.string({
            char: 't',
            description: chalk.green.bold('Type of channel.'),
            options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
            required: false,
        }),
        config: flags.boolean({
            char: 'c',
            description: chalk.green.bold('Return channel config values only.'),
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(ChannelGet);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json'),
                this
            );
            const channel = await client.queryChannels(
                { id: flags.id, type: flags.type },
                { last_message_at: -1 },
                {
                    subscribe: false,
                }
            );

            const payload = cardinal.highlight(
                stringify(
                    flags.config ? channel[0].data.config : channel[0].data,
                    {
                        maxLength: 100,
                    },
                    {
                        linenos: true,
                    }
                )
            );

            this.log(payload);
            this.exit(0);
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ChannelGet.description = 'Get a channel.';
