const { Command, flags } = require('@oclif/command');
const fs = require('fs');

const { auth } = require('../../../../utils/auth');

class PushSetApn extends Command {
	async run() {
		const { flags } = this.parse(PushSetApn);

		try {
			const client = await auth(this);

			const settings = await client.updateAppSettings({
				apn_config: {
					p12_cert: fs.readFileSync(flags.p12_cert) || '',
					pem_cert: fs.readFileSync(flags.pem_cert, 'utf-8') || '',
					auth_key: flags.auth_key || '',
					key_id: flags.key_id || '',
					team_id: flags.team_id || '',
					topic: flags.bundle_id || '',
					notification_template: flags.notification_template,
				},
			});

			if (flags.json) {
				this.log(JSON.stringify(settings));
				this.exit();
			}

			this.log('Push notifications have been enabled with APN.');
			this.exit(0);
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

PushSetApn.flags = {
	auth_key: flags.string({
		char: 'a',
		description: 'Private auth key.',
		required: false,
	}),
	key_id: flags.string({
		char: 'k',
		description: 'Key ID.',
		required: false,
	}),
	bundle_id: flags.string({
		char: 'b',
		description: 'Bundle identifier (e.g. com.apple.test).',
		required: false,
	}),
	team_id: flags.string({
		char: 't',
		description: 'Team ID.',
		required: false,
	}),
	pem_cert: flags.string({
		char: 'p',
		description: 'Absolute path to .pem RSA key.',
		required: false,
	}),
	p12_cert: flags.string({
		char: 'b',
		description: 'Absolute path to .p12 file.',
		required: false,
	}),
	notification_template: flags.string({
		char: 'n',
		description: 'JSON template for notifications.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.PushSetApn = PushSetApn;
