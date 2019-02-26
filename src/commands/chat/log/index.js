const { Command, flags } = require('@oclif/command');
const stringify = require('json-stringify-pretty-compact');
const prompt = require('enquirer');
const cardinal = require('cardinal');
const moment = require('moment');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

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

class Log extends Command {
    async run() {
        const { flags } = this.parse(Log);

        try {
            const client = await auth(this);

            if (!flags.event) {
                const res = await prompt({
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

                flags.event = res.event;
            }

            await client.updateUser({
                id: '*',
                role: 'admin',
            });

            await client.setUser({
                id: '*',
                status: 'invisible',
            });

            const channel = client.channel(flags.type, flags.channel);
            await channel.watch();

            this.log(`Logging real-time events for ${flags.event}...}`);

            const time = 'dddd, MMMM Do YYYY [at] h:mm:ss A';

            if (flags.event === 'all') {
                channel.on(event => {
                    const timestamp = chalk.yellow.bold(
                        moment(event.created_at).format(time)
                    );

                    const payload = `${timestamp}: ${chalk.bold(
                        event.user.name || event.user.id
                    )} (${chalk.bold(
                        event.user.role
                    )}) performed event ${chalk.bold(
                        event.type
                    )} in channel ${chalk.bold(flags.channel)}.`;

                    this.log(payload);
                });
            } else {
                channel.on(flags.event, event => {
                    const timestamp = chalk.yellow.bold(
                        moment(event.created_at).format(time)
                    );

                    const payload = cardinal.highlight(
                        stringify(event, { maxLength: 100 }),
                        { linenos: true }
                    );

                    this.log(`${timestamp}:`, '\n\n', payload, '\n\n');
                });
            }
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

Log.flags = {
    channel: flags.string({
        char: 'c',
        description: 'The channel ID you wish to log.',
        required: false,
    }),
    type: flags.string({
        char: 't',
        description: 'The type of channel.',
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: true,
    }),
    event: flags.string({
        char: 'e',
        description: 'The type of event you want to listen on.',
        options: events,
        required: false,
    }),
};

module.exports.Log = Log;
