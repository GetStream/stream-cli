const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class ModerateMute extends Command {
    async run() {
        const { flags } = this.parse(ModerateMute);

        try {
            const client = await auth(this);
            const mute = await client.muteUser(flags.user);

            if (flags.json) {
                this.log(mute);
                this.exit(0);
            }

            this.log(`The message ${chalk.bold(flags.user)} has been flagged!`);
            this.exit(0);
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

ModerateMute.flags = {
    user: flags.string({
        char: 'u',
        description: 'The ID of the user to mute.',
        required: true,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.ModerateMute = ModerateMute;
