const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class UserUpdate extends Command {
	async run() {
		const { flags } = this.parse(UserUpdate);

		try {
			if (!flags.user || !flags.name) {
				const res = await prompt([
					{
						type: 'input',
						name: 'id',
						message: 'What is the unique identifier for the user?',
						required: true,
					},
					{
						type: 'input',
						name: 'name',
						message: 'What is the name of the user?',
						required: true,
					},
					{
						type: 'input',
						name: 'image',
						message: 'What is the URL of the users image?',
						required: false,
					},
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await chatAuth(this);
			const token = client.createToken(flags.id);

			const payload = {
				id: flags.id,
				name: flags.name,
			};

			if (flags.image) {
				payload.image = flags.image;
			}

			await client.setUser(payload, token);

			if (flags.json) {
				this.log(JSON.stringify(payload));
				this.exit();
			}

			this.log(`The user ${flags.name} (${flags.id}) has been set.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

UserUpdate.flags = {
	id: flags.string({
		char: 'i',
		description: 'The unique identifier for the user.',
		required: false,
	}),
	name: flags.string({
		char: 'n',
		description: 'Name of the user.',
		required: false,
	}),
	image: flags.string({
		char: 'm',
		description: 'URL to the image of the user.',
		required: false,
	}),
};

UserUpdate.description = 'Updates a user.';

module.exports.UserUpdate = UserUpdate;
