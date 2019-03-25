const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const treeify = require('treeify');

const { auth } = require('../../../utils/auth');

class ChannelGet extends Command {
    async run() {
        const { flags } = this.parse(ChannelGet);

        try {
            if (!flags.channel || !flags.type) {
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
                ]);

                for (const key in res) {
                    if (res.hasOwnProperty(key)) {
                        flags[key] = res[key];
                    }
                }
            }

            const client = await auth(this);

            const channel = await client.queryChannels(
                { id: flags.channel, type: flags.type },
                { last_message_at: -1 },
                {
                    subscribe: false,
                }
            );

            if (flags.json) {
                this.log(JSON.stringify(channel[0].data));
                this.exit(0);
            }

            delete channel[0].data.config['commands'];
            delete channel[0].data.config['created_at'];
            delete channel[0].data.config['updated_at'];

            const tree = treeify.asTree(channel[0].data, true, false);

            this.log(tree);
            this.exit(0);
        } catch (error) {
            this.error(error.message || 'A Stream CLI error has occurred.', {
                exit: 1,
            });
        }
    }
}

ChannelGet.flags = {
    channel: flags.string({
        char: 'c',
        description: 'The channel ID you wish to retrieve.',
        required: false,
    }),
    type: flags.string({
        char: 't',
        description: 'Type of channel.',
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.ChannelGet = ChannelGet;
