const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');

const { auth } = require('../../../utils/auth');

class ModerateBan extends Command {
    async run() {
        const { flags } = this.parse(ModerateBan);

        try {
            const client = await auth(this);
            const ban = await client.banUser(flags.user, {
                timeout: Number(flags.timeout),
                reason: flags.reason,
            });

            if (flags.json) {
                this.log(JSON.stringify(ban));
                this.exit(0);
            }

            this.log(`The user ${chalk.bold(flags.user)} has been banned.`);
            this.exit(0);
        } catch (error) {
            this.error(error || 'A Stream CLI error has occurred.', {
                exit: 1,
            });
        }
    }
}

ModerateBan.flags = {
    user: flags.string({
        char: 'u',
        description: 'The ID of the offending user.',
        exclusive: ['message'],
        required: true,
    }),
    reason: flags.string({
        char: 'r',
        description: 'A reason for adding a timeout.',
        required: true,
    }),
    timeout: flags.string({
        char: 't',
        description: 'Duration of timeout in minutes.',
        default: '60',
        required: true,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.ModerateBan = ModerateBan;
