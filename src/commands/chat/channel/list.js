const { Command, flags } = require('@oclif/command');
const treeify = require('treeify');
const moment = require('moment');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class ChannelList extends Command {
    async run() {
        const { flags } = this.parse(ChannelList);

        try {
            const client = await auth(this);

            const channels = await client.queryChannels(
                {},
                { last_message_at: -1 },
                {
                    watch: false,
                    state: false,
                    subscribe: false,
                }
            );

            if (flags.json) {
                for (const channel of channels) {
                    this.log(channel, '\n');
                }

                this.exit(0);
            }

            for (const channel of channels) {
                delete channel.data.config['commands'];
                delete channel.data.config['created_at'];
                delete channel.data.config['updated_at'];

                const tree = treeify.asTree(channel.data, true, false);

                this.log(tree);
            }

            this.exit(0);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

ChannelList.flags = {
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.ChannelList = ChannelList;
