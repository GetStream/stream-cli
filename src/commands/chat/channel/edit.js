const { Command, flags } = require('@oclif/command');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');
const uuid = require('uuid/v4');

const { auth } = require('../../../utils/auth');

class ChannelEdit extends Command {
    async run() {
        const { flags } = this.parse(ChannelEdit);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json'),
                this
            );
            const channel = await client.channel(flags.type, flags.id);

            let payload = {
                name: flags.name,
                updated_by: {
                    id: uuid(),
                    name: 'CLI',
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

            this.log(
                `The channel ${flags.id} has been modified!`,
                emoji.get('rocket')
            );
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ChannelEdit.flags = {
    id: flags.string({
        char: 'i',
        description: chalk.blue.bold('The ID of the channel you wish to edit.'),
        required: true,
    }),
    type: flags.string({
        char: 't',
        description: chalk.blue.bold('Type of channel.'),
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: true,
    }),
    name: flags.string({
        char: 'n',
        description: chalk.blue.bold('Name of the channel room.'),
        required: true,
    }),
    url: flags.string({
        char: 'u',
        description: chalk.blue.bold('URL to the channel image.'),
        required: false,
    }),
    reason: flags.string({
        char: 'r',
        description: chalk.blue.bold('Reason for changing channel.'),
        required: true,
    }),
    members: flags.string({
        char: 'm',
        description: chalk.blue.bold('Comma separated list of members.'),
        required: false,
    }),
    data: flags.string({
        char: 'd',
        description: chalk.blue.bold('Additional data as JSON.'),
        required: false,
    }),
};

module.exports.ChannelEdit = ChannelEdit;
