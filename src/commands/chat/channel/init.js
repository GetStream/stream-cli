const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');
const path = require('path');
const uuid = require('uuid/v4');

const { auth } = require('../../../utils/auth');
const { credentials } = require('../../../utils/config');

class ChannelInit extends Command {
    async run() {
        const { flags } = this.parse(ChannelInit);

        try {
            const { name, email } = await credentials(this);
            const client = await auth(this);

            let payload = {
                name: flags.name,
                created_by: {
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

            const channel = await client.channel(
                flags.type,
                flags.channel,
                payload
            );

            const create = await channel.create();

            if (flags.json) {
                this.log(create);
                this.exit(0);
            }

            this.log(
                `The channel ${chalk.bold(flags.name)} has been initialized.`
            );
            this.exit(0);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

ChannelInit.flags = {
    channel: flags.string({
        char: 'c',
        description: 'A unique ID for the channel you wish to create.',
        default: uuid(),
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
    image: flags.string({
        char: 'i',
        description: 'URL to channel image.',
        required: false,
    }),
    members: flags.string({
        char: 'm',
        description: 'Comma separated list of members to add to the channel.',
        required: false,
    }),
    data: flags.string({
        char: 'd',
        description: 'Additional data as a JSON.',
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.ChannelInit = ChannelInit;
