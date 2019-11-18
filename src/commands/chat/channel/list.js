const { Command, flags } = require('@oclif/command');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class ChannelList extends Command {
	async run() {
		const { flags } = this.parse(ChannelList);

		try {
			const client = await chatAuth(this);
			const channels = await client.queryChannels(
				{},
				{ last_message_at: -1 },
				{
					watch: false,
					state: false,
					subscribe: false,
					limit: parseInt(flags.limit, 10) || 10,
					offset: parseInt(flags.offset, 10) || 0,
				}
			);

			const arr = [];
			for (const c of channels) {
				arr.push(c.data);
			}

			this.log(JSON.stringify(arr));

			this.exit();
		} catch (error) {
			await this.config.runHook('telemetry', {
				ctx: this,
				error,
			});
		}
	}
}

ChannelList.flags = {
	limit: flags.string({
		char: 'l',
		description: 'Channel list limit.',
		required: true,
	}),
	offset: flags.string({
		char: 'o',
		description: 'Channel list offset.',
		required: true,
	}),
};

ChannelList.description = 'Lists all channels.';

module.exports.ChannelList = ChannelList;
