const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class ChannelQuery extends Command {
    async run() {
        const { flags } = this.parse(ChannelQuery);

        try {
            const client = await auth(this);

            const filter = flags.filters ? JSON.parse(flags.filters) : {};
            const sort = flags.sort ? JSON.parse(flags.sort) : {};

            const channels = await client.queryChannels(filter, sort, {
                subscribe: false,
            });

            this.log(channels[0].data);
            this.exit(0);
        } catch (err) {
            this.error(err || 'A CLI error has occurred.', { exit: 1 });
        }
    }
}

ChannelQuery.flags = {
    id: flags.string({
        char: 'i',
        description: chalk.blue.bold('The channel ID you wish to query.'),
        required: true,
    }),
    type: flags.string({
        char: 't',
        description: chalk.blue.bold('Type of channel.'),
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: false,
    }),
    filter: flags.string({
        char: 'f',
        description: chalk.blue.bold('Filters to apply to the query.'),
        required: false,
    }),
    sort: flags.string({
        char: 's',
        description: chalk.blue.bold('Sort to apply to the query.'),
        required: false,
    }),
};

module.exports.ChannelQuery = ChannelQuery;
