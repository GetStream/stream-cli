const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');
const { credentials } = require('../../../utils/config');

class ChannelUpdate extends Command {
    async run() {
        const { flags } = this.parse(ChannelUpdate);

        try {
            const { name, email } = await credentials(this);

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

            const update = await channel.update(payload, {
                name: flags.name,
                text: flags.reason,
            });

            if (flags.json) {
                this.log(update);
                this.exit(0);
            }

            this.log(`The channel ${chalk.bold(flags.id)} has been modified.`);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

ChannelUpdate.flags = {
    id: flags.string({
        char: 'i',
        description: 'The ID of the channel you wish to update.',
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
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.ChannelUpdate = ChannelUpdate;
