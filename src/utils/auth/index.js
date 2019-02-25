const { StreamChat } = require('stream-chat');

const { credentials } = require('../../utils/config');

async function auth(config) {
    try {
        const { apiKey, apiSecret } = await credentials(config);
        if (!apiKey || !apiSecret) {
            throw new Error('Missing configuration file...');
        }

        const client = new StreamChat(apiKey, apiSecret);

        return client;
    } catch (err) {
        throw new Error(err);
    }
}

module.exports.auth = auth;
