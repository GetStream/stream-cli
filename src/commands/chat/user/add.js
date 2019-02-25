const { Command, flags } = require('@oclif/command');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class UserAdd extends Command {
    async run() {
        const { flags } = this.parse(UserAdd);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

            const channel = await client.channel(flags.type, flags.id);
            await channel.addModerators(flags.moderators.split(','));

            this.log(
                `${flags.moderators} have been added as moderators to channel ${
                    flags.type
                }:${flags.id}`,
                emoji.get('rocket')
            );
            this.exit(0);
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

UserAdd.flags = {
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
            'Comma separated list of moderators to add.'
        ),
        required: true,
    }),
};

module.exports.UserAdd = UserAdd;
