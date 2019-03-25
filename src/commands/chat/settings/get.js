const { Command, flags } = require('@oclif/command');

const { auth } = require('../../../utils/auth');

class SettingsGet extends Command {
	async run() {
		try {
			this.exit(0);
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

SettingsGet.flags = {
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.SettingsGet = SettingsGet;
