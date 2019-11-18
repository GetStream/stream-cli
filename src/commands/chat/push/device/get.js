const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const Table = require('cli-table');

const { chatAuth } = require('../../../../utils/auth/chat-auth');

class DeviceList extends Command {
	async run() {
		const { flags } = this.parse(DeviceList);

		try {
			if (!flags.user_id) {
				const result = await prompt([
					{
						type: 'input',
						name: 'user_id',
						hint: 'user-123',
						message: 'What is the User ID?',
						required: true,
					},
				]);

				for (const key in result) {
					if (result.hasOwnProperty(key)) {
						flags[key] = result[key];
					}
				}
			}

			const client = await chatAuth(this);

			const response = await client.getDevices(flags.user_id || '');

			if (response.devices && response.devices.length !== 0) {
				const table = new Table({
					head: ['Device ID', 'Push provider'],
				});
				for (const device of response.devices) {
					table.push([device.id, device.push_provider]);
				}
				this.log(table.toString());
			} else {
				this.log('User has no devices');
			}
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

DeviceList.flags = {
	user_id: flags.string({
		char: 'u',
		description: 'User ID',
		required: false,
	}),
};

DeviceList.description = 'Gets all devices registered for push.';

module.exports.DeviceList = DeviceList;
