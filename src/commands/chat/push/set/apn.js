const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');
const fs = require('fs');

const { auth } = require('../../../../utils/auth');

class PushSetApn extends Command {
	async run() {
		const { flags } = this.parse(PushSetApn);

		try {
			if (!flags.p12_cert && !flags.pem_cert && !flags.auth_key) {
				const type = await prompt([
					{
						type: 'select',
						name: 'apn',
						message: 'What type of APN configuration will you use?',
						required: true,
						choices: [
							{ message: '.p12', value: 'p12' },
							{ message: '.pem', value: 'pem' },
							{ message: '.p8', value: 'p8' },
						],
					},
				]);

				let res = {};

				if (type.apn === 'p12') {
					res = await prompt([
						{
							type: 'input',
							name: 'p12_cert',
							hint: '/Users/username/Desktop/cert.p12',
							message:
								'What is the absolute path to your .p12 file?',
							required: true,
						},
						{
							type: 'input',
							name: 'bundle_id',
							hint: 'com.apple.test',
							message: 'What is your bundle ID?',
							required: true,
						},
						{
							type: 'input',
							name: 'notification_template',
							hint: 'Omit for Stream default',
							message: `What JSON notification template would you like to use?`,
							required: false,
						},
					]);
				}

				if (type.apn === 'pem') {
					res = await prompt([
						{
							type: 'input',
							name: 'pem_cert',
							hint: '/Users/username/Desktop/cert.pem',
							message:
								'What is the absolute path to your .pem file?',
							required: true,
						},
						{
							type: 'input',
							name: 'bundle_id',
							hint: 'com.apple.test',
							message: 'What is your bundle ID?',
							required: true,
						},
						{
							type: 'input',
							name: 'notification_template',
							hint: 'Omit for Stream default',
							message: `What JSON notification template would you like to use?`,
							required: false,
						},
					]);
				}

				if (type.apn === 'p8') {
					res = await prompt([
						{
							type: 'input',
							name: 'auth_key',
							hint: '/Users/username/Desktop/key.p8',
							message: 'What is the filepath to your .p8 file?',
							required: true,
						},
						{
							type: 'input',
							name: 'key_id',
							message: 'What is your key ID?',
							required: true,
						},
						{
							type: 'input',
							name: 'bundle_id',
							hint: 'com.apple.test',
							message: 'What is your bundle ID?',
							required: true,
						},
						{
							type: 'input',
							name: 'team_id',
							message: 'What is your team ID?',
							required: true,
						},
						{
							type: 'input',
							name: 'notification_template',
							hint: 'Omit for Stream default',
							message: `What JSON notification template would you like to use?`,
							required: false,
						},
					]);
				}

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await auth(this);

			const payload = {
				apn_config: {
					p12_cert: fs.existsSync(flags.p12_cert)
						? fs.readFileSync(flags.p12_cert)
						: '',
					pem_cert: fs.existsSync(flags.pem_cert)
						? fs.readFileSync(flags.pem_cert, 'utf-8')
						: '',
					auth_key: fs.existsSync(flags.auth_key)
						? fs.readFileSync(flags.auth_key, 'utf-8')
						: '',
					key_id: flags.key_id || '',
					team_id: flags.team_id || '',
					topic: flags.bundle_id || '',
				},
			};

			if (flags.notification_template) {
				payload.apn_config.notification_template =
					flags.notification_template;
			}

			await client.updateAppSettings(payload);

			if (flags.json) {
				const settings = await client.getAppSettings();

				this.log(JSON.stringify(settings.app.push_notifications));
				this.exit();
			}

			this.log(
				`Push notifications have been enabled with ${chalk.bold(
					'APN'
				)}.`
			);
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
		description: 'Absolute path to .p8 auth key.',
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
