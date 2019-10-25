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

		const settings = environment === 'production' ? false : true;

		await chatClient.updateAppSettings({
			disable_auth_checks: settings,
			disable_permissions_checks: settings,
		});

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
