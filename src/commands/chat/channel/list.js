const { Command, flags } = require('@oclif/command');
const treeify = require('treeify');

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
                const arr = [];

                for (const c of channels) {
                    arr.push(c.data);
                }

                this.log(JSON.stringify(arr));
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
        } catch (error) {
            this.error(error.message || 'A Stream CLI error has occurred.', {
                exit: 1,
            });
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
