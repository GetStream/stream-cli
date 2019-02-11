import { Command, flags } from '@oclif/command';
import { StreamChat } from 'stream-chat';
import stringify from 'json-stringify-pretty-compact';
import { prompt } from 'enquirer';
import cardinal from 'cardinal';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import fs from 'fs-extra';

import { authError } from '../../utils/error';
import { credentials } from '../../utils/config';

const events = [
    'all',
    'user.status.changed',
    'user.watching.start',
    'user.watching.stop',
    'user.updated',
    'typing.start',
    'typing.stop',
    'message.new',
    'message.updated',
    'message.deleted',
    'message.seen',
    'message.reaction',
    'member.added',
    'member.removed',
    'channel.updated',
    'health.check',
    'connection.changed',
    'connection.recovered',
];

export class Log extends Command {
    static flags = {
        type: flags.string({
            char: 't',
            description: 'Channel type to tail.',
            required: true,
        }),
        channel: flags.string({
            char: 'c',
            description: 'Channel name to tail.',
            required: true,
        }),
        event: flags.string({
            char: 'e',
            description: 'Event type to tail.',
            options: events,
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(Log);
        const config = path.join(this.config.configDir, 'config.json');

        try {
            const { apiKey, apiSecret } = await credentials(config);
            if (!apiKey || !apiSecret) return authError();

            const client = new StreamChat(apiKey, apiSecret);

            if (!flags.event) {
                const question = await prompt({
                    type: 'autocomplete',
                    name: 'event',
                    message: 'What event would you like to filter on?',
                    limit: 17,
                    suggest(input, choices) {
                        return choices.filter(choice =>
                            choice.message.startsWith(input)
                        );
                    },
                    choices: events,
                });

                flags.event = question.event;
            }

            await client.updateUser({
                id: '*',
                role: 'admin',
            });

            const result = await client.setUser({
                id: '*',
                status: 'invisible',
            });

            const channel = client.channel(flags.type, flags.channel);

            await channel.watch();

            console.log(
                chalk.green.bold(
                    `Logging real-time events for ${flags.event}... ${emoji.get(
                        'rocket'
                    )}`
                )
            );

            const time = 'dddd, MMMM Do YYYY [at] h:mm:ss A';

            if (flags.event === 'all') {
                let timestamp = chalk.yellow.bold(
                    moment(event.channel.last_message_at).format(time)
                );

                channel.on(event => {
                    let payload = `${timestamp}: ${chalk.green.bold(
                        event.user.id
                    )} (${chalk.green.bold(
                        event.user.role
                    )}) performed event ${chalk.green.bold(
                        event.type
                    )} in channel ${chalk.green.bold(flags.channel)}.`;

                    console.log(payload);
                });
            } else {
                channel.on(flags.event, event => {
                    let timestamp = chalk.yellow.bold(
                        moment(event.channel.last_message_at).format(time)
                    );

                    let payload = cardinal.highlight(
                        stringify(event, { maxLength: 100 }),
                        { linenos: true }
                    );

                    console.log(`${timestamp}:`, '\n\n', payload, '\n\n');
                });
            }
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

Log.description = 'watch events in real-time coming from a team channel';
