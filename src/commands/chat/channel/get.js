const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const treeify = require('treeify');
const numeral = require('numeral');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class ChannelGet extends Command {
    async run() {
        const { flags } = this.parse(ChannelGet);
        const client = await auth(this);

        try {
            if (!flags.channel || !flags.type) {
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
                ]);

                for (const key in res) {
                    if (res.hasOwnProperty(key)) {
                        flags[key] = res[key];
                    }
                }
            }

            const channel = await client.queryChannels(
                { id: flags.channel, type: flags.type },
                { last_message_at: -1 },
                {
                    subscribe: false,
                }
            );

            if (flags.raw) {
                this.log(channel[0]);
                this.exit(0);
            }

            delete channel[0].data.config['commands'];
            delete channel[0].data.config['created_at'];
            delete channel[0].data.config['updated_at'];

            const tree = treeify.asTree(channel[0].data, true, false);

            this.log(tree);
            this.exit(0);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

ChannelGet.flags = {
    channel: flags.string({
        char: 'c',
        description: chalk.blue.bold('The channel ID you wish to get.'),
        required: false,
    }),
    type: flags.string({
        char: 't',
        description: chalk.blue.bold('Type of channel.'),
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: false,
    }),
    raw: flags.string({
        char: 'r',
        description: chalk.blue.bold(
            'A raw object containing all channel data.'
        ),
        required: false,
    }),
};

module.exports.ChannelGet = ChannelGet;
