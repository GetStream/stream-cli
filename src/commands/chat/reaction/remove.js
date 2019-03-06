const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');

const { auth } = require('../../../utils/auth');

class ReactionRemove extends Command {
    async run() {
        const { flags } = this.parse(ReactionRemove);

        try {
            if (!flags.json) {
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
                        type: 'input',
                        name: 'message',
                        message: `What is the unique identifier for the message?`,
                        required: true,
                    },
                    {
                        type: 'input',
                        name: 'reaction',
                        message: `What is the unique identifier for the reaction?`,
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
            const channel = client.channel(flags.type, flags.channel);

            const reaction = await channel.deleteReaction(
                flags.message,
                flags.reaction
            );

            if (flags.json) {
                this.log(JSON.stringify(reaction));
                this.exit(0);
            }

            this.log('The reaction has been removed.');
            this.exit(0);
        } catch (error) {
            this.error(error || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

ReactionRemove.flags = {
    channel: flags.string({
        char: 'c',
        description: 'The unique identifier for the channel.',
        required: false,
    }),
    type: flags.string({
        char: 't',
        description: 'The type of channel.',
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: false,
    }),
    message: flags.string({
        char: 'c',
        description: 'The unique identifier for the message.',
        required: false,
    }),
    reaction: flags.string({
        char: 'r',
        description: 'The unique identifier for the reaction.',
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.ReactionRemove = ReactionRemove;
