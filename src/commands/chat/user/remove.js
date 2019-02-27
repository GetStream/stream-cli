const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class UserRemove extends Command {
    async run() {
        const { flags } = this.parse(UserRemove);

        try {
            const client = await auth(this);

            const channel = await client.channel(flags.type, flags.channel);
            const remove = await channel.demoteModerators(
                flags.moderators.split(',')
            );

            if (flags.json) {
                this.log(remove);
                this.exit(0);
            }

            this.log(
                `${chalk.bold(
                    flags.moderators
                )} have been removed as moderators from the ${chalk.bold(
                    flags.type
                )} channel ${chalk.bold(flags.channel)}.`
            );
            this.exit(0);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

UserRemove.flags = {
    channel: flags.string({
        char: 'c',
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
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.UserRemove = UserRemove;
