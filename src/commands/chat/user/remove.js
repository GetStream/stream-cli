const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class UserRemove extends Command {
    async run() {
        const { flags } = this.parse(UserRemove);

        try {
            const client = await auth(this);

            const channel = await client.channel(flags.type, flags.id);
            await channel.demoteModerators(flags.moderators.split(','));

            this.log(
                `${chalk.bold(
                    flags.moderators
                )} have been removed as moderators from the ${chalk.bold(
                    flags.type
                )} channel ${chalk.bold(flags.id)}.`
            );
            this.exit(0);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

UserRemove.flags = {
    id: flags.string({
        char: 'i',
        description: 'Channel name.',
        required: true,
    }),
    type: flags.string({
        char: 't',
        description: 'Channel type.',
        required: true,
    }),
    moderators: flags.string({
        char: 'm',
        description: 'Comma separated list of moderators to remove.',
        required: true,
    }),
};

module.exports.UserRemove = UserRemove;
