const { StreamChat } = require('stream-chat');
const chalk = require('chalk');

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

		return chatClient;
	} catch (error) {
		ctx.error(
			`Authentication required. Use the command ${chalk.green.bold(
				'stream config:set'
			)} to authenticate.`,
			{
				exit: 1,
			}
		);
	}
}

module.exports.chatAuth = chatAuth;
