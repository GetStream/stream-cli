const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');

const { auth } = require('../../../utils/auth');

class MessageRemove extends Command {
    async run() {
        const { flags } = this.parse(MessageRemove);

        try {
            if (!flags.message || !flags.json) {
                const res = await prompt([
                    {
                        type: 'input',
                        name: 'message',
                        message: `What is the unique identifier for the message?`,
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
            const remove = await client.deleteMessage(flags.message);

            if (flags.json) {
                this.log(JSON.stringify(remove));
                this.exit(0);
            }

            this.log(
                `The message ${chalk.bold(flags.message)} has been removed.`
            );
            this.exit(0);
        } catch (error) {
            this.error(error || 'A Stream CLI error has occurred.', {
                exit: 1,
            });
        }
    }
}

MessageRemove.flags = {
    message: flags.string({
        char: 'message',
        description: 'The ID of the message you would like to remove.',
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.MessageRemove = MessageRemove;
