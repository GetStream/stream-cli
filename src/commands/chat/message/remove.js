const { Command, flags } = require('@oclif/command');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class MessageRemove extends Command {
    async run() {
        const { flags } = this.parse(MessageRemove);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

            await client.deleteMessage(flags.id);

            this.log(
                `The message ${flags.id} has been removed!`,
                emoji.get('wastebasket')
            );
            this.exit(0);
        } catch (err) {
            this.error(err || 'A CLI error has occurred.', { exit: 1 });
        }
    }
}

MessageRemove.flags = {
    channel: flags.string({
        char: 'c',
        description: chalk.blue.bold(
            'The channel ID that you would like to remove.'
        ),
        required: true,
    }),
};

module.exports.MessageRemove = MessageRemove;
