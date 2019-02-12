import { Command, flags } from '@oclif/command';
import stringify from 'json-stringify-pretty-compact';
import { prompt } from 'enquirer';
import cardinal from 'cardinal';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';
import uuid from 'uuid/v4';

import { auth } from '../../utils/auth';
import { apiError } from '../../utils/error';

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
        id: flags.string({
            char: 'i',
            description: chalk.green.bold('Channel ID.'),
            default: uuid(),
            required: false,
        }),
        type: flags.string({
            char: 't',
            description: chalk.green.bold('Channel type.'),
            required: true,
        }),
        event: flags.string({
            char: 'e',
            description: chalk.green.bold('Event type.'),
            options: events,
            required: false,
        }),
    };

    async run() {
        const { flags } = this.parse(Log);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

            if (!flags.event) {
                const question = await prompt({
                    type: 'autocomplete',
                    name: 'event',
                    message: 'What event would you like to filter on?',
                    limit: events.length,
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

            const channel = client.channel(flags.type, flags.id);

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
                channel.on(event => {
                    let timestamp = chalk.yellow.bold(
                        moment(event.created_at).format(time)
                    );

                    let payload = `${timestamp}: ${chalk.green.bold(
                        event.user.id
                    )} (${chalk.green.bold(
                        event.user.role
                    )}) performed event ${chalk.green.bold(
                        event.type
                    )} in channel ${chalk.green.bold(flags.id)}.`;

                    console.info(payload);
                });
            } else {
                channel.on(flags.event, event => {
                    let timestamp = chalk.yellow.bold(
                        moment(event.created_at).format(time)
                    );

                    let payload = cardinal.highlight(
                        stringify(event, { maxLength: 100 }),
                        { linenos: true }
                    );

                    console.info(`${timestamp}:`, '\n\n', payload, '\n\n');
                });
            }
        } catch (err) {
            apiError(err);
        }
    }
}

Log.description = 'watch events in real-time coming from a team channel';
