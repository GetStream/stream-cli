const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class UserCreate extends Command {
    async run() {
        const { flags } = this.parse(UserCreate);

        try {
            if (!flags.type || !flags.moderators || !flags.channel) {
                const res = await prompt([
                    {
                        type: 'input',
                        name: 'channel',
                        hint: 'The name of the channel',
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
                        type: 'input',
                        name: 'moderators',
                        message: 'Who would you like to add as a moderator?',
                        hint: 'e.g. Thierry, Tommaso, Nick (Comma Separated)',
                        required: true,
                    },
                ]);

                for (const key in res) {
                    if (res.hasOwnProperty(key)) {
                        flags[key] = res[key];
                    }
                }
            }

            const client = await auth(this);
            const channel = await client.channel(flags.type, flags.channel);
            const create = await channel.addModerators(
                flags.moderators.split(',')
            );

            if (flags.json) {
                this.log(create);
                this.exit(0);
            }

            this.log(
                `${chalk.bold(
                    flags.moderators
                )} have been added as moderators to channel ${chalk.bold(
                    flags.type
                )}:${chalk.bold(flags.channel)}`
            );
            this.exit(0);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

UserCreate.flags = {
    channel: flags.string({
        char: 'c',
        description: 'Channel identifier.',
        required: false,
    }),
    type: flags.string({
        char: 't',
        description: 'The type of channel.',
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: false,
    }),
    moderators: flags.string({
        char: 'm',
        description: 'Comma separated list of moderators.',
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.UserCreate = UserCreate;
