const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');

const { auth } = require('../../../utils/auth');

class MessageQuery extends Command {
    async run() {
        const { flags } = this.parse(MessageQuery);

        try {
            const client = await auth(this);

            if (flags.json) {
                this.log(add);
                this.exit(0);
            }

            this.log(message);
            this.exit();
        } catch (err) {
            this.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
        }
    }
}

MessageQuery.flags = {
    json: flags.boolean({
        char: 'j',
        description:
            'Output results in JSON. When not specified, returns output in a human friendly format.',
        required: false,
    }),
};

module.exports.MessageQuery = MessageQuery;
