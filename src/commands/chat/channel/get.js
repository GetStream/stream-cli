import { Command, flags } from '@oclif/command';
import Table from 'cli-table';
import numeral from 'numeral';
import chalk from 'chalk';
import path from 'path';

import { auth } from '../../../utils/auth';

export class ChannelGet extends Command {
    static flags = {
        id: flags.string({
            char: 'i',
            description: chalk.blue.bold('The channel ID you wish to get.'),
            required: true,
        }),
        type: flags.string({
            char: 't',
            description: chalk.blue.bold('Type of channel.'),
            options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
            required: true,
        }),
    };

    async run() {
        const { flags } = this.parse(ChannelGet);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json'),
                this
            );

            const channel = await client.queryChannels(
                { id: flags.id, type: flags.type },
                { last_message_at: -1 },
                {
                    subscribe: false,
                }
            );

            const table = new Table();

            const data = channel[0].data;
            const config = data.config;

            table.push(
                { CID: data.cid },
                { Name: data.name },
                { Type: data.type },
                {
                    Owner: `${data.created_by.name} (${data.created_by.role})`,
                },
                {
                    Roles: data.channel_roles.length
                        ? data.channel_roles.join(', ')
                        : chalk.red('No roles defined'),
                },
                {
                    Members: data.members.length
                        ? data.members.join(', ')
                        : chalk.red('No active members'),
                },
                {
                    Automod:
                        config.automod === 'enabled'
                            ? chalk.green('enabled')
                            : chalk.red('disabled'),
                },
                {
                    Commands: config.commands
                        .map(
                            command =>
                                `${command.name} (${chalk.green('enabled')})`
                        )
                        .join(', '),
                },
                {
                    Mutes:
                        config.mutes === true
                            ? chalk.green('enabled')
                            : chalk.red('disabled'),
                },
                {
                    Reactions:
                        config.reactions === true
                            ? chalk.green('enabled')
                            : chalk.red('disabled'),
                },
                {
                    Replies:
                        config.replies === true
                            ? chalk.green('enabled')
                            : chalk.red('disabled'),
                },
                {
                    Search:
                        config.search === true
                            ? chalk.green('enabled')
                            : chalk.red('disabled'),
                },
                {
                    Events: [
                        config.connect_events === true
                            ? `connect (${chalk.green('enabled')})`
                            : '',
                        config.seen_events === true
                            ? `seen (${chalk.green('enabled')})`
                            : '',
                        config.typing_events === true
                            ? `typing (${chalk.green('enabled')})`
                            : '',
                    ].join(', '),
                },
                { 'Message Retention': config.message_retention },
                {
                    'Max Message Length': `${numeral(
                        config.max_message_length
                    ).format('0,0')} characters`,
                }
            );

            this.log(table.toString());
            this.exit(0);
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ChannelGet.description = 'Get a specific channel by ID.';
