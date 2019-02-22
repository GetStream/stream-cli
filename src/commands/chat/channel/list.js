const { Command, flags } = require('@oclif/command');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class ChannelList extends Command {
    async run() {
        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json'),
                this
            );

            const channels = await client.queryChannels(
                {},
                { last_message_at: -1 },
                {
                    subscribe: false,
                }
            );

            if (channels.length) {
                channels.map(channel =>
                    this.log(
                        chalk.blue(
                            `The channel ${chalk.bold(
                                channel.id
                            )} of type ${chalk.bold(
                                channel.type
                            )} with the CID of ${chalk.bold(
                                channel.cid
                            )} has ${chalk.bold(
                                channel.data.members.length
                            )} members.`
                        )
                    )
                );

                this.exit(0);
            } else {
                this.warn(
                    `Your application does not have any channels.`,
                    emoji.get('pensive')
                );

                this.exit(0);
            }
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

module.exports.ChannelList = ChannelList;
