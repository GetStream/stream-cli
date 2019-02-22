const { Command, flags } = require('@oclif/command');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

export class MessageRemove extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold(
                'The channel ID that you would like to remove.'
            ),
            required: true,
        }),
    };

    async run() {
        const { flags } = this.parse(MessageRemove);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json'),
                this
            );

            await client.deleteMessage(flags.id);

            this.log(
                `The message ${flags.id} has been removed!`,
                emoji.get('wastebasket')
            );
            this.exit(0);
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

MessageRemove.description = 'Remove messages from a channel.';
