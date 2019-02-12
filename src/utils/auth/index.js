import { StreamChat } from 'stream-chat';

import { authError } from '../../utils/error';
import { credentials } from '../../utils/config';

export async function auth(config) {
    try {
        const { apiKey, apiSecret } = await credentials(config);
        if (!apiKey || !apiSecret) return authError();

        const client = new StreamChat(apiKey, apiSecret);

        return client;
    } catch (err) {
        authError(err);
    }
}
