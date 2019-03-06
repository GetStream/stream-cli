const { Command, flags } = require('@oclif/command');
const stringify = require('json-stringify-pretty-compact');
const { prompt } = require('enquirer');
const cardinal = require('cardinal');
const moment = require('moment');
const chalk = require('chalk');

const { auth } = require('../../../utils/auth');

class Log extends Command {
    async run() {
        const { flags } = this.parse(Log);

        try {
            if (!flags.channel || !flags.type || !flags.event) {
                const res = await prompt([
                    {
                        type: 'input',
                        name: 'channel',
                        message: `What is the unique identifier for the channel?`,
                        required: true,
                    },
                    {
                        type: 'select',
                        name: 'type',
                        message: 'What type of channel is this?',
                        required: true,
                        choices: [
                            { message: 'Livestream', value: 'livestream' },
                            { message: 'Messaging', value: 'messaging' },
                            { message: 'Gaming', value: 'gaming' },
                            { message: 'Commerce', value: 'commerce' },
                            { message: 'Team', value: 'team' },
                        ],
                    },
                    {
                        type: 'select',
                        name: 'event',
                        message: 'What event would you like to filter on?',
                        required: true,
                        choices: [
                            {
                                message: 'All Events - JSON',
                                value: 'all',
                            },
                            {
                                message: 'User Status - Changed',
                                value: 'user.status.changed',
                            },
                            {
                                message: 'User Watching - Start',
                                value: 'user.watching.start',
                            },
                            {
                                message: 'User Watching - Stop',
                                value: 'user.watching.stop',
                            },
                            {
                                message: 'User Updated',
                                value: 'user.updated',
                            },
                            {
                                message: 'Typing - Start',
                                value: 'typing.start',
                            },
                            {
                                message: 'Typing - Stop',
                                value: 'typing.stop',
                            },
                            {
                                message: 'Message - New',
                                value: 'message.new',
                            },
                            {
                                message: 'Message - Updated',
                                value: 'message.updated',
                            },
                            {
                                message: 'Message - Deleted',
                                value: 'message.deleted',
                            },
                            {
                                message: 'Message - Seen',
                                value: 'message.seen',
                            },
                            {
                                message: 'Message - Reaction',
                                value: 'message.reaction',
                            },
                            {
                                message: 'Member - Added',
                                value: 'member.added',
                            },
                            {
                                message: 'Member - Removed',
                                value: 'member.removed',
                            },
                            {
                                message: 'Channel - Updated',
                                value: 'channel.updated',
                            },
                            {
                                message: 'Health - Check',
                                value: 'health.check',
                            },
                            {
                                message: 'Connection - Changed',
                                value: 'connection.changed',
                            },
                            {
                                message: 'Connection - Recovered',
                                value: 'connection.recovered',
                            },
                        ],
                    },
                ]);

                for (const key in res) {
                    if (res.hasOwnProperty(key)) {
                        flags[key] = res[key];
                    }
                }
            }

            const client = await auth(this);

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

            if (flags.json) {
                channel.on(event => {
                    this.log(JSON.stringify(event));
                });
            } else if (flags.event === 'all') {
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
        } catch (error) {
            this.error(error || 'A Stream CLI error has occurred.', {
                exit: 1,
            });
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
        required: false,
    }),
    event: flags.string({
        char: 'e',
        description: 'The type of event you want to listen on.',
        options: [
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
        ],
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.Log = Log;
