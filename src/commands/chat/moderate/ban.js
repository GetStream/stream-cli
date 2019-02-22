const { Command, flags } = require('@oclif/command');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

export class ModerateBan extends Command {
    static flags = {
        user: flags.string({
            char: 'u',
            description: chalk.blue.bold('The ID of the offending user.'),
            exclusive: ['message'],
            required: true,
        }),
        reason: flags.string({
            char: 'r',
            description: chalk.blue.bold('A reason for adding a timeout.'),
            required: true,
        }),
        timeout: flags.string({
            char: 't',
            description: chalk.blue.bold('Duration of timeout in minutes.'),
            default: '60',
            required: true,
        }),
    };

    async run() {
        const { flags } = this.parse(ModerateBan);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json'),
                this
            );

            await client.banUser(flags.user, {
                timeout: Number(flags.timeout),
                reason: flags.reason,
            });

            this.log(
                `The user ${flags.user} has been banned!`,
                emoji.get('banned')
            );
            this.exit(0);
        } catch (err) {
            this.error(err, { exit: 1 });
        }
    }
}

ModerateBan.description =
    'Ban users from a channel indefinitely or based on a per minute timeout.';
