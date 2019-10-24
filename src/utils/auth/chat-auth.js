const { StreamChat } = require('stream-chat');

const { credentials } = require('../../utils/config');

async function chatAuth(ctx) {
	try {
		const {
			apiKey,
			apiSecret,
			apiBaseUrl,
			environment,
		} = await credentials(ctx);

		const chatClient = new StreamChat(apiKey, apiSecret);
		chatClient.setBaseURL(apiBaseUrl);

		await chatClient.updateAppSettings({
			disable_auth_checks: environment === 'production' ? false : true,
			disable_permissions_checks:
				environment === 'production' ? false : true,
		});

		return chatClient;
	} catch (error) {
		ctx.error(error || 'A Stream authentication error has occurred.', {
			exit: 1,
		});
	}
}

module.exports.chatAuth = chatAuth;
