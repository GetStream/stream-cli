const StreamChat = require('stream-chat');

const { credentials } = require('../../utils/config');

async function auth(config, _this) {
    try {
        const { apiKey, apiSecret } = await credentials(config);
        if (!apiKey || !apiSecret) {
            return _this.error(`Missing config...`, { exit: 1 });
        }

        const client = new StreamChat(apiKey, apiSecret);

        return client;
    } catch (err) {
        _this.error(err, { exit: 1 });
    }
}

module.exports.auth = auth;
