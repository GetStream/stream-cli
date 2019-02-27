const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class MessageRemove extends Command {
    async run() {
        const { flags } = this.parse(MessageRemove);

        try {
            const client = await auth(this);
            const remove = await client.deleteMessage(flags.id);

            if (flags.json) {
                this.log(remove);
                this.exit(0);
            }

            this.log(`The message ${chalk.bold(flags.id)} has been removed.`);
            this.exit(0);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

MessageRemove.flags = {
    channel: flags.string({
        char: 'c',
        description: 'The channel ID you are targeting.',
        required: true,
    }),
    message: flags.string({
        char: 'message',
        description: 'The ID of the message you would like to remove.',
        required: true,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.MessageRemove = MessageRemove;
