const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');

const { auth } = require('../../../utils/auth');

class ChannelInvite extends Command {
    async run() {
        const { flags } = this.parse(ChannelInvite);

        try {
            const client = await auth(this);
            const channel = await client.channel(flags.type, flags.id);
            const invite = await client.updateUsers(flags.users.split(','));

            if (flags.json) {
                this.log('JSON');
                this.exit(0);
            }

            this.log(`The member ${chalk.bold(flags.name)} has been invited.`);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

ChannelInvite.flags = {
    channel: flags.string({
        char: 'c',
        description: 'The ID of the channel you wish to invite a user to.',
        required: false,
    }),
    type: flags.string({
        char: 't',
        description: 'Type of channel.',
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: false,
    }),
    members: flags.string({
        char: 'm',
        description: 'Comma separated list of members to invite.',
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.ChannelInvite = ChannelInvite;
