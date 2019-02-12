import { Command, flags } from '@oclif/command';
import stringify from 'json-stringify-pretty-compact';
import cardinal from 'cardinal';
import emoji from 'node-emoji';
import moment from 'moment';
import chalk from 'chalk';
import path from 'path';

import { auth } from '../../utils/auth';
import { exit } from '../../utils/response';
import { apiError } from '../../utils/error';

export class ChannelList extends Command {
    async run() {
        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
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

            console.log(`${ts}\n`);

            if (channels.length) {
                channels.map(channel => {
                    console.log(
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

                process.exit(0);
            } else {
                exit(`Your application does not have any channels.`, {
                    emoji: 'pensive',
                });
            }
        } catch (err) {
            apiError(err);
        }
    }
}

ChannelList.description = 'Get a channel.';
