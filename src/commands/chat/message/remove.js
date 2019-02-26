const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class MessageRemove extends Command {
    async run() {
        const { flags } = this.parse(MessageRemove);

        try {
            const client = await auth(this);

            await client.deleteMessage(flags.id);

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
};

module.exports.MessageRemove = MessageRemove;
