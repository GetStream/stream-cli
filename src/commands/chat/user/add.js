const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class UserAdd extends Command {
    async run() {
        const { flags } = this.parse(UserAdd);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

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

            const channel = await client.channel(flags.type, flags.channel);
            await channel.addModerators(flags.moderators.split(','));

            this.log(
                `${flags.moderators} have been added as moderators to channel ${
                    flags.type
                }:${flags.channel}`,
                emoji.get('rocket')
            );
            this.exit(0);
        } catch (err) {
            this.error(err || 'A CLI error has occurred.', { exit: 1 });
        }
    }
}

UserAdd.flags = {
    channel: flags.string({
        char: 'c',
        description: chalk.blue.bold('Channel identifier.'),
        required: false,
    }),
    type: flags.string({
        char: 't',
        description: chalk.blue.bold('The type of channel.'),
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: false,
    }),
    moderators: flags.string({
        char: 'm',
        description: chalk.blue.bold('Comma separated list of moderators.'),
        required: false,
    }),
};

module.exports.UserAdd = UserAdd;
