import { StreamChat } from 'stream-chat';
import chalk from 'chalk';

import { credentials } from '../../utils/config';

async function chatAuth(ctx) {
	try {
		const { apiKey, apiSecret, timeout } = await credentials(
			ctx
		);

		const chatClient = new StreamChat(apiKey, apiSecret, { timeout });
		chatClient.setBaseURL("https://chat.stream-io-api.com");

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
