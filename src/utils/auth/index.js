import { StreamChat } from 'stream-chat';

import { credentials } from '../../utils/config';

export async function auth(config, _this) {
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
