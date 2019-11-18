const { Command, flags } = require('@oclif/command');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class ChannelQuery extends Command {
	async run() {
		const { flags } = this.parse(ChannelQuery);

		try {
			const client = await chatAuth(this);

			const filter = flags.filter ? JSON.parse(flags.filter) : {};
			const sort = flags.sort ? JSON.parse(flags.sort) : {};

			const channel = await client.queryChannels(filter, sort, {
				subscribe: false,
			});

			if (flags.json) {
				this.log(JSON.stringify(channel[0].data));
				this.exit();
			}

			this.log(channel[0].data);
			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

ChannelQuery.flags = {
	channel: flags.string({
		char: 'c',
		description: 'The unique identifier for the channel you want to query.',
		required: false,
	}),
	type: flags.string({
		char: 't',
		description: 'Type of channel.',
		required: false,
	}),
	filter: flags.string({
		char: 'f',
		description: 'Filters to apply to the query.',
		required: false,
	}),
	sort: flags.string({
		char: 's',
		description: 'Sort to apply to the query.',
		required: false,
	}),
	json: flags.boolean({
		char: 'j',
		description:
			'Output results in JSON. When not specified, returns output in a human friendly format.',
		required: false,
	}),
};

ChannelQuery.description = 'Queries all channels.';

module.exports.ChannelQuery = ChannelQuery;
