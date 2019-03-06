const { Command, flags } = require('@oclif/command');
const chalk = require('chalk');

const { auth } = require('../../../utils/auth');

class ModerateFlag extends Command {
    async run() {
        const { flags } = this.parse(ModerateFlag);

        try {
            const client = await auth(this);

            if (flags.user) {
                await client.flagUser(flags.user);

                this.log(`User ${chalk.bold(flags.user)} has been flagged!`);
                this.exit(0);
            } else if (flags.message) {
                const flag = await client.flagMessage(flags.message);

                if (flags.json) {
                    this.log(JSON.stringify(flag));
                    this.exit(0);
                }

                this.log(
                    `Message ${chalk.bold(flags.message)} has been flagged!`
                );
                this.exit(0);
            } else {
                this.warn(
                    `Please pass a valid command. Use the command ${chalk.bold(
                        'chat:moderate:flag --help'
                    )} for more information.`
                );
                this.exit(0);
            }
        } catch (error) {
            this.error(error || 'A Stream CLI error has occurred.', {
                exit: 1,
            });
        }
    }
}

ModerateFlag.flags = {
    user: flags.string({
        char: 'u',
        description: 'The ID of the offending user.',
        exclusive: ['message'],
        required: false,
    }),
    message: flags.string({
        char: 'm',
        description: 'The ID of the message you want to flag.',
        exclusive: ['user'],
        required: false,
    }),
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.ModerateFlag = ModerateFlag;
