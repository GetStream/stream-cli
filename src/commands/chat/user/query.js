const { Command, flags } = require('@oclif/command');
const { prompt } = require('enquirer');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class UserQuery extends Command {
	async run() {
		const { flags } = this.parse(UserQuery);

		try {
			if (!flags.query) {
				const query = await prompt([
					{
						type: 'input',
						name: 'query',
						message: 'What is the query you would like to perform?',
						hint: 'optional',
						required: false,
					},
				]);

				flags.query = query;
			}

			if (!flags.sort) {
				const sort = await prompt([
					{
						type: 'select',
						name: 'sort',
						message: 'What user key would you like to sort on?',
						choices: [
							{
								message: 'ID',
								value: 'id',
							},
							{
								message: 'Last Active',
								value: 'last_active',
							},
							{
								message: 'None',
								value: 'none',
							},
						],
						required: false,
					},
				]);

				flags.sort = sort;
			}

			if (!flags.limit) {
				const limit = await prompt([
					{
						type: 'input',
						name: 'query',
						message: 'How many records would you like to display?',
						required: false,
					},
				]);

				flags.limit = limit;
			}

			if (!flags.offset) {
				const offset = await prompt([
					{
						type: 'input',
						name: 'query',
						message: 'How many users would you like to skip?',
						required: false,
					},
				]);

				flags.offset = offset;
			}

			const s = {};
			if (flags.sort.sort !== 'none') {
				switch (flags.sort) {
					case 'id':
						s.id = -1;
						break;
					case 'last_active':
						s.last_active = -1;
						break;
					default:
						s.id = -1;
				}
			}

			let q = {};
			if (flags.query.length) {
				q = JSON.parse(flags.query);
			}

			const client = await chatAuth(this);
			const { users } = await client.queryUsers(q, s, {
				limit: parseInt(flags.limit.query, 10) || 25,
				offset: parseInt(flags.offset.query, 10) || 0,
			});

			if (flags.json) {
				this.log(JSON.stringify(users));
				this.exit();
			}

			this.log(users);
			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

UserQuery.flags = {
	query: flags.string({
		char: 'q',
		description: 'The query you would like to perform.',
		required: false,
	}),
	sort: flags.string({
		char: 's',
		description: 'Display the current status of the user.',
		required: false,
	}),
	limit: flags.string({
		char: 'l',
		description: 'The limit to apply to the query.',
		required: false,
	}),
	offset: flags.string({
		char: 'o',
		description: 'The offset to apply to the query.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

module.exports.UserQuery = UserQuery;
