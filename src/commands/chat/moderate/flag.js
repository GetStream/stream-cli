const { Command, flags } = require('@oclif/command');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');

const { auth } = require('../../../utils/auth');

class ModerateFlag extends Command {
    async run() {
        const { flags } = this.parse(ModerateFlag);

        try {
            const client = await auth(
                path.join(this.config.configDir, 'config.json')
            );

            if (flags.user) {
                await client.flagUser(flags.user);

                this.log(
                    `The user ${flags.user} has been flagged!`,
                    emoji.get('bangbang')
                );
                this.exit(0);
            } else if (flags.message) {
                await client.flagMessage(flags.message);

                this.log(
                    `The message ${flags.user} has been flagged!`,
                    emoji.get('bangbang')
                );
                this.exit(0);
            } else {
                this.warn(
                    `Please pass a valid command. Use the command ${chalk.blue.bold(
                        'moderate:flag --help'
                    )} for more information.`
                );
                this.exit(0);
            }
        } catch (err) {
            this.error(err || 'A CLI error has occurred.', { exit: 1 });
        }
    }
}

ModerateFlag.flags = {
    user: flags.string({
        char: 'u',
        description: chalk.blue.bold('The ID of the offending user.'),
        exclusive: ['message'],
        required: false,
    }),
    message: flags.string({
        char: 'm',
        description: chalk.blue.bold('The ID of the message you want to flag.'),
        exclusive: ['user'],
        required: false,
    }),
};

module.exports.ModerateFlag = ModerateFlag;
