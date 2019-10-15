const { StreamChat } = require('stream-chat');

const { credentials } = require('../../utils/config');

async function chatAuth(ctx) {
	try {
		const { apiKey, apiSecret, apiBaseUrl } = await credentials(ctx);

		const chatClient = new StreamChat(apiKey, apiSecret);
		chatClient.setBaseURL(apiBaseUrl);

		return chatClient;
	} catch (error) {
		ctx.error(error || 'A Stream authentication error has occurred.', {
			exit: 1,
		});
	}
}

module.exports.chatAuth = chatAuth;
