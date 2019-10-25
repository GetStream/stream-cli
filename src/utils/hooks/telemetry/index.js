const md5 = require('md5');
const Rollbar = require('rollbar');

const { credentials } = require('../../config');

module.exports = async ({ ctx, error }) => {
	const {
		name,
		email,
		apiKey,
		apiSecret,
		apiBaseUrl,
		environment,
		telemetry,
	} = await credentials(ctx);

	if (telemetry === 'true') {
		const rollbar = new Rollbar({
			accessToken: '4ba9cceb33fe4543b7a114fc354b870c',
			captureUncaught: true,
			captureUnhandledRejections: true,
			environment: 'production',
		});

		rollbar.error(error, {
			person: {
				id: md5(email),
				username: email,
				name,
				email,
			},
			api: {
				key: apiKey,
				secret: apiSecret,
				url: apiBaseUrl,
			},
			environment,
		});
	}

	ctx.error(`A Stream CLI error has occurred: ${error.message}`, {
		exit: 1,
	});
};
