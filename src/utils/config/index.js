const emoji = require('node-emoji');
const chalk = require('chalk');
const fs = require('fs-extra');

async function credentials(config) {
    try {
        if (!(await fs.pathExists(config))) {
            await fs.outputJson(config, {
                apiKey: '',
                apiSecret: '',
            });
        }

        const { apiKey, apiSecret } = await fs.readJson(config);

        if (!apiKey.length || !apiSecret.length) {
            console.warn(
                `Credentials not found. Run ${chalk.bold(
                    'stream config:set'
                )} to generate a new configuration file.`,
                emoji.get('warning')
            );

            process.exit(0);
        }

        return { apiKey, apiSecret };
    } catch (err) {
        console.warn(err);
        process.exit(1);
    }
}

module.exports.credentials = credentials;
