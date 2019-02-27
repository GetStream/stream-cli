const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');
const uuid = require('uuid/v4');

const { auth } = require('../../../utils/auth');
const { credentials } = require('../../../utils/config');

class MessageCreate extends Command {
    async run() {
        const { flags } = this.parse(MessageCreate);

        try {
            const { name } = await credentials(this);

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

            const client = await auth(this);
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

            const add = await channel.sendMessage(payload);

            if (flags.json) {
                this.log(add);
                this.exit(0);
            }

            const message = `Message ${chalk.bold(
                flags.message
            )} has been sent to the ${chalk.bold(
                flags.channel
            )} channel by user ${chalk.bold(flags.user)}.`;

            this.log(message);
            this.exit();
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

MessageCreate.flags = {
    user: flags.string({
        char: 'u',
        description: 'The ID of the user sending the message.',
        required: false,
    }),
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
    message: flags.string({
        char: 'm',
        description: 'The message you would like to send as plaintext.',
        required: false,
    }),
    attachments: flags.string({
        char: 'a',
        description:
            'A JSON payload of attachments to send along with a message.',
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.MessageCreate = MessageCreate;
