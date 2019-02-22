const emoji = require('node-emoji');
const chalk = require('chalk');
const fs = require('fs-extra');

export async function credentials(config, _this) {
    try {
        if (!(await fs.pathExists(config))) {
            await fs.outputJson(config, {
                apiKey: '',
                apiSecret: '',
            });
        }

        const { apiKey, apiSecret } = await fs.readJson(config);

        if (!apiKey.length || !apiSecret.length) {
            _this.log(
                chalk.red(
                    `Credentials not found. Run ${chalk.bold(
                        'stream config:set'
                    )} to generate a configuration file. ${emoji.get(
                        'warning'
                    )}`
                )
            );

            _this.exit(0);
        }

        return { apiKey, apiSecret };
    } catch (err) {
        _this.error(err);
        _this.exit(1);
    }
}
