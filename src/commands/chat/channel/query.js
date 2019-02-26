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
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

ChannelQuery.flags = {
    channel: flags.string({
        char: 'c',
        description: 'The channel ID you wish to query.',
        required: true,
    }),
    type: flags.string({
        char: 't',
        description: 'Type of channel.',
        options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
        required: false,
    }),
    filter: flags.string({
        char: 'f',
        description: 'Filters to apply to the query.',
        required: false,
    }),
    sort: flags.string({
        char: 's',
        description: 'Sort to apply to the query.',
        required: false,
    }),
};

module.exports.ChannelQuery = ChannelQuery;
