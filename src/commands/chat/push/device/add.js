const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');

const { chatAuth } = require('../../../../utils/auth/chat-auth');

class DeviceAdd extends Command {
	async run() {
		const { flags } = this.parse(DeviceAdd);

		try {
			if (!flags.user_id || !flags.provider || !flags.device_id) {
				const result = await prompt([
					{
						type: 'input',
						name: 'user_id',
						hint: 'user-123',
						message: 'What is the User ID?',
						required: true,
					},
					{
						type: 'input',
						name: 'device_id',
						hint: `device-123`,
						message: 'What is the Device ID?',
						required: true,
					},
					{
						type: 'select',
						name: 'provider',
						message: 'What is the push provider?',
						required: true,
						choices: [
							{ message: 'APN', value: 'apn' },
							{ message: 'Firebase', value: 'firebase' },
						],
					},
				]);

				for (const key in result) {
					if (result.hasOwnProperty(key)) {
						flags[key] = result[key];
					}
				}
			}

			const client = await chatAuth(this);

			await client.addDevice(
				flags.device_id || '',
				flags.provider || '',
				flags.user_id || ''
			);

			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

DeviceAdd.flags = {
	user_id: flags.string({
		char: 'u',
		description: 'User ID',
		required: false,
	}),
	device_id: flags.string({
		char: 'd',
		description: 'Device id or token.',
		required: false,
	}),
	provider: flags.string({
		char: 'p',
		description: 'Push provider',
		required: false,
	}),
};

module.exports.DeviceAdd = DeviceAdd;
