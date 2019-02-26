const { StreamChat } = require('stream-chat');

const { credentials } = require('../../utils/config');

async function auth(ctx) {
    try {
        const { apiKey, apiSecret } = await credentials(ctx);

        const client = new StreamChat(apiKey, apiSecret);

        return client;
    } catch (err) {
        ctx.error(err || 'A Stream authentication error has occurred.', {
            exit: 1,
        });
    }
}

module.exports.auth = auth;
