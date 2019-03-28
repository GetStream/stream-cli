const { StreamChat } = require('stream-chat');

const { credentials } = require('../../utils/config');

async function chatAuth(ctx) {
	try {
		const { apiKey, apiSecret } = await credentials(ctx);

		return new StreamChat(apiKey, apiSecret);
	} catch (error) {
		ctx.error(error || 'A Stream authentication error has occurred.', {
			exit: 1,
		});
	}
}

module.exports.chatAuth = chatAuth;
