const { Command, flags } = require('@oclif/command');
const Table = require('cli-table');
const numeral = require('numeral');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class ChannelGet extends Command {
    async run() {
        const { flags } = this.parse(ChannelGet);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
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
                { [`${chalk.green.bold('CID')}`]: data.cid },
                { [`${chalk.green.bold('Name')}`]: data.name },
                { [`${chalk.green.bold('Type')}`]: data.type },
                {
                    [`${chalk.green.bold('Owner')}`]: `${
                        data.created_by.name
                    } (${data.created_by.role})`,
                },
                {
                    [`${chalk.green.bold('Roles')}`]: data.channel_roles.length
                        ? data.channel_roles.join(', ')
                        : chalk.red('No roles defined'),
                },
                {
                    [`${chalk.green.bold('Members')}`]: data.members.length
                        ? data.members.join(', ')
                        : chalk.red('No active members'),
                },
                {
                    [`${chalk.green.bold('Automod')}`]:
                        config.automod === 'enabled'
                            ? chalk.green('enabled')
                            : chalk.red('disabled'),
                },
                {
                    [`${chalk.green.bold('Commands')}`]: config.commands
                        .map(
                            command =>
                                `${command.name} (${chalk.green('enabled')})`
                        )
                        .join(', '),
                },
                {
                    [`${chalk.green.bold('Mutes')}`]:
                        config.mutes === true
                            ? chalk.green('enabled')
                            : chalk.red('disabled'),
                },
                {
                    [`${chalk.green.bold('Reactions')}`]:
                        config.reactions === true
                            ? chalk.green('enabled')
                            : chalk.red('disabled'),
                },
                {
                    [`${chalk.green.bold('Replies')}`]:
                        config.replies === true
                            ? chalk.green('enabled')
                            : chalk.red('disabled'),
                },
                {
                    [`${chalk.green.bold('Search')}`]:
                        config.search === true
                            ? chalk.green('enabled')
                            : chalk.red('disabled'),
                },
                {
                    [`${chalk.green.bold('Events')}`]: [
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
                {
                    [`${chalk.green.bold(
                        'Message Retention'
                    )}`]: config.message_retention,
                },
                {
                    [`${chalk.green.bold('Max Message Length')}`]: `${numeral(
                        config.max_message_length
                    ).format('0,0')} characters`,
                }
            );

            this.log(table.toString());
            this.exit(0);
        } catch (err) {
            this.error(err || 'A CLI error has occurred.', { exit: 1 });
        }
    }
}

ChannelGet.flags = {
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

module.exports.ChannelGet = ChannelGet;
