const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');
const uuid = require('uuid/v4');
const path = require('path');

const { auth } = require('../../../utils/auth');
const { credentials } = require('../../../utils/config');

class MessageSend extends Command {
    async run() {
        const { name } = await credentials(this);

        const { flags } = this.parse(MessageSend);

        try {
            const client = await auth(this);

            if (!flags.user || !flags.channel || !flags.message || flags.type) {
                const res = await prompt([
                    {
                        type: 'input',
                        name: 'user',
                        message: `What is the unique identifier for the user sending this message?`,
                        default: name || uuid(),
                        required: true,
                    },
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
                        name: 'message',
                        hint: 'Hello World!',
                        message: 'What is the message you would like to send?',
                        required: true,
                    },
                ]);

                for (const key in res) {
                    if (res.hasOwnProperty(key)) {
                        flags[key] = res[key];
                    }
                }
            }

            await client.updateUser({
                id: flags.user,
                role: 'admin',
            });

            await client.setUser({ id: flags.user, status: 'invisible' });
            const channel = client.channel(flags.type, flags.channel);

            const payload = {
                text: flags.message,
            };

            if (flags.attachments) {
                payload.attachments = JSON.parse(flags.attachments);
            }

            await channel.sendMessage(payload);

            const message = `Message ${chalk.bold(
                flags.message
            )} has been sent to the ${chalk.bold(
                flags.channel
            )} channel by user ${chalk.bold(flags.user)}!`;

            this.log(message);
            this.exit();
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

MessageSend.flags = {
    user: flags.string({
        char: 'u',
        description: chalk.blue.bold('The ID of the user sending the message.'),
        required: false,
    }),
    type: flags.string({
        char: 't',
        description: chalk.blue.bold('The type of channel.'),
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: false,
    }),
    channel: flags.string({
        char: 'c',
        description: chalk.blue.bold(
            'The ID of the channel that you would like to send a message to.'
        ),
        required: false,
    }),
    message: flags.string({
        char: 'm',
        description: chalk.blue.bold(
            'The message you would like to send as plaintext.'
        ),
        required: false,
    }),
    attachments: flags.string({
        char: 'a',
        description: chalk.blue.bold(
            'A JSON payload of attachments to send along with a message.'
        ),
        required: false,
    }),
};

module.exports.MessageSend = MessageSend;
