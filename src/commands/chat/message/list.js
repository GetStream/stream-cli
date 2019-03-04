const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { auth } = require('../../../utils/auth');

class MessageList extends Command {
    async run() {
        const { flags } = this.parse(MessageList);

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
                ]);

                for (const key in res) {
                    if (res.hasOwnProperty(key)) {
                        flags[key] = res[key];
                    }
                }
            }

            const client = await auth(this);
            const channel = client.channel(flags.type, flags.channel);

            const messages = await client.queryChannels(
                {},
                { last_message_at: -1 }
            );

            if (flags.json) {
                if (!messages.length) {
                    this.log(messages);
                    this.exit(0);
                }

                this.log(messages[0].state.messages);
                this.exit(0);
            }

            const data = messages[0].state.messages;

            if (!data.length) {
                this.log('No messages available.');
                this.exit();
            }

            data.map(message => {
                this.log(
                    `${chalk.bold.green(message.id)} (${message.created_at}): ${
                        message.text
                    }`
                );
            });

            this.exit();
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

MessageList.flags = {
    type: flags.string({
        char: 't',
        description: 'The type of channel.',
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: false,
    }),
    channel: flags.string({
        char: 'c',
        description:
            'The ID of the channel that you would like to send a message to.',
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.MessageList = MessageList;
