const { Command } = require('@oclif/command');

const { chatAuth } = require('../../../utils/auth/chat-auth');

class ChannelList extends Command {
	async run() {
		try {
			const client = await chatAuth(this);
			const channels = await client.queryChannels(
				{},
				{ last_message_at: -1 },
				{
					watch: false,
					state: false,
					subscribe: false,
				}
			);

			const arr = [];
			for (const c of channels) {
				arr.push(c.data);
			}

			this.log(JSON.stringify(arr));

			this.exit();
		} catch (error) {
			this.error(error || 'A Stream CLI error has occurred.', {
				exit: 1,
			});
		}
	}
}

module.exports.ChannelList = ChannelList;
