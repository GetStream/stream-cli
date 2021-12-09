import md5 from 'md5';
import Rollbar from 'rollbar';

import { credentials } from '../../config';

module.exports = async ({ ctx, error }) => {
	const { name, email, apiKey, apiSecret, environment, telemetry } = await credentials(ctx);

	if (telemetry === 'true') {
		const rollbar = new Rollbar({
			accessToken: '4ba9cceb33fe4543b7a114fc354b870c',
			captureUncaught: true,
			captureUnhandledRejections: true,
			environment: 'production'
		});

		rollbar.error(error, {
			person: {
				id: md5(email),
				username: email,
				name,
				email
			},
			api: {
				key: apiKey,
				secret: apiSecret,
				url: "https://chat.stream-io-api.com",
			},
			environment
		});
	}

	if (typeof error === 'object' && (error.code !== 'EEXIT' || error.code === undefined)) {
		ctx.error(error.message, {
			exit: 1
		});
	}
};
