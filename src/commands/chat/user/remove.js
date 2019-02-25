const { Command, flags } = require('@oclif/command');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class UserRemove extends Command {
    async run() {
        const { flags } = this.parse(UserRemove);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

            const channel = await client.channel(flags.type, flags.id);
            await channel.demoteModerators(flags.moderators.split(','));

            this.log(
                `${flags.moderators} have been removed as moderators from the ${
                    flags.type
                } channel ${flags.id}`,
                emoji.get('warning')
            );
            this.exit(0);
        } catch (err) {
            this.error(err || 'A CLI error has occurred.', { exit: 1 });
        }
    }
}

UserRemove.flags = {
    id: flags.string({
        char: 'i',
        description: chalk.blue.bold('Channel name.'),
        required: true,
    }),
    type: flags.string({
        char: 't',
        description: chalk.blue.bold('Channel type.'),
        required: true,
    }),
    moderators: flags.string({
        char: 'm',
        description: chalk.blue.bold(
            'Comma separated list of moderators to remove.'
        ),
        required: true,
    }),
};

module.exports.UserRemove = UserRemove;
