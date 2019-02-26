const { Command } = require('@oclif/command');
const emoji = require('node-emoji');
const chalk = require('chalk');
const path = require('path');
const fs = require('fs-extra');

async function credentials(ctx) {
    const config = path.join(ctx.config.configDir, 'config.json');

    try {
        if (!(await fs.pathExists(config))) {
            await fs.outputJson(config, {
                name: '',
                email: '',
                apiKey: '',
                apiSecret: '',
            });
        }

        const { name, email, apiKey, apiSecret } = await fs.readJson(config);

        if (!name || !email || !apiKey || !apiSecret) {
            console.warn(
                `Credentials not found. Run the command ${chalk.bold(
                    'stream config:set'
                )} to generate a new Stream configuration file.`,
                emoji.get('warning')
            );

            ctx.exit(0);
        }

        return { name, email, apiKey, apiSecret };
    } catch (err) {
        ctx.error(err || 'A Stream CLI error has occurred.', { exit: 1 });
    }
}

module.exports.credentials = credentials;
