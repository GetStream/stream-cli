const { StreamChat } = require('stream-chat');

const { credentials } = require('../../utils/config');

async function auth(ctx) {
    try {
        const { apiKey, apiSecret } = await credentials(ctx);

        return new StreamChat(apiKey, apiSecret);
    } catch (err) {
        ctx.error(err || 'A Stream authentication error has occurred.', {
            exit: 1,
        });
    }
}

module.exports.auth = auth;
