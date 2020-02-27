const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');
const chalk = require('chalk');
const uuid = require('uuid/v4');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class UserCreate extends Command {
	async run() {
		const { flags } = this.parse(UserCreate);

		try {
			if (!flags.user || !flags.role) {
				const res = await prompt([
					{
						type: 'input',
						name: 'user',
						message: 'What is the unique identifier for the user?',
						default: uuid(),
						required: true
					},
					{
						type: 'select',
						name: 'role',
						message: 'What role would you like assign to the user?',
						required: true,
						choices: [
							{
								message: 'Admin',
								value: 'admin'
							},
							{
								message: 'User',
								value: 'user'
							}
						]
					}
				]);

				for (const key in res) {
					if (res.hasOwnProperty(key)) {
						flags[key] = res[key];
					}
				}
			}

			const client = await chatAuth(this);
			const create = await client.updateUser({
				id: flags.user,
				role: flags.role
			});

			if (flags.json) {
				this.log(JSON.stringify(create));
				this.exit();
			}

			this.log(`The user ${chalk.bold(flags.user)} (${flags.role}) has been created.`);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error
			});
		}
	}
}

UserCreate.flags = {
	user: flags.string({
		char: 'u',
		description: 'Comma separated list of users to add.',
		required: false
	}),
	role: flags.string({
		char: 'r',
		description: 'The role to assign to the user.',
		options: [ 'admin', 'user' ],
		required: false
	}),
	json: flags.boolean({
		char: 'j',
		description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false
	})
};

UserCreate.description = 'Creates a new user.';

module.exports.UserCreate = UserCreate;
