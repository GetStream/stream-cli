const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');
const { credentials } = require('../../../utils/config');

class ChannelEdit extends Command {
    async run() {
        const { name, email } = await credentials(this);

        const { flags } = this.parse(ChannelEdit);

        try {
            const client = await auth(this);
            const channel = await client.channel(flags.type, flags.id);

            let payload = {
                name: flags.name,
                updated_by: {
                    id: email,
                    name,
                },
            };
            if (flags.image) payload.image = flags.image;
            if (flags.members) payload.members = flags.members.split(',');

            if (flags.data) {
                const parsed = JSON.parse(flags.data);
                payload = Object.assign({}, payload, parsed);
            }

            await channel.update(payload, {
                name: flags.name,
                text: flags.reason,
            });

            this.log(`The channel ${chalk.bold(flags.id)} has been modified.`);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

ChannelEdit.flags = {
    id: flags.string({
        char: 'i',
        description: 'The ID of the channel you wish to edit.',
        required: true,
    }),
    type: flags.string({
        char: 't',
        description: 'Type of channel.',
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: true,
    }),
    name: flags.string({
        char: 'n',
        description: 'Name of the channel room.',
        required: true,
    }),
    url: flags.string({
        char: 'u',
        description: 'URL to the channel image.',
        required: false,
    }),
    reason: flags.string({
        char: 'r',
        description: 'Reason for changing channel.',
        required: true,
    }),
    members: flags.string({
        char: 'm',
        description: 'Comma separated list of members.',
        required: false,
    }),
    data: flags.string({
        char: 'd',
        description: 'Additional data as JSON.',
        required: false,
    }),
};

module.exports.ChannelEdit = ChannelEdit;
