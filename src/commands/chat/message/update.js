const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');
const uuid = require('uuid/v4');

const { auth } = require('../../../utils/auth');

class MessageUpdate extends Command {
    async run() {
        const { flags } = this.parse(MessageUpdate);

        try {
            if (!flags.message || !flags.text) {
                const res = await prompt([
                    {
                        type: 'input',
                        name: 'message',
                        message: `What is the unique identifier for the message?`,
                        default: uuid(),
                        required: true,
                    },
                    {
                        type: 'input',
                        name: 'text',
                        message: 'What is the updated message?',
                        required: true,
                    },
                ]);

                for (const key in res) {
                    if (res.hasOwnProperty(key)) {
                        flags[key] = res[key];
                    }
                }
            }

            const payload = {
                id: flags.message,
                text: flags.text,
            };

            if (flags.attachments) {
                payload.attachments = JSON.parse(flags.attachments);
            }

            const client = await auth(this);
            const update = await client.updateMessage(payload);

            if (flags.json) {
                this.log(JSON.stringify(update));
                this.exit(0);
            }

            this.log(
                `Message ${chalk.bold(flags.message.id)} has been updated.`
            );
            this.exit();
        } catch (error) {
            this.error(error || 'A Stream CLI error has occurred.', {
                exit: 1,
            });
        }
    }
}

MessageUpdate.flags = {
    message: flags.string({
        char: 'm',
        description: 'The unique identifier for the message.',
        required: false,
    }),
    text: flags.string({
        char: 't',
        description: 'The message you would like to send as text.',
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

module.exports.MessageUpdate = MessageUpdate;
