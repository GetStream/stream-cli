import { Command } from '@oclif/command';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';

import { auth } from '../../utils/auth';

export class ChannelList extends Command {
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

            const ts = chalk.yellow.bold(
                moment().format('dddd, MMMM Do YYYY [at] h:mm:ss A')
            );

            this.log(`${ts}\n`);

            if (channels.length) {
                return channels.map(channel => {
                    this.log(
                        chalk.blue(
                            `The Channel ${chalk.bold(
                                channel.id
                            )} of type ${chalk.bold(
                                channel.type
                            )} with the CID of ${chalk.bold(
                                channel.cid
                            )} has ${chalk.bold(
                                channel.data.members.length
                            )} members.`
                        )
                    );
                });

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

ChannelList.description = 'Lists all channels.';
