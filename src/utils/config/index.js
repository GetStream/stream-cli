import emoji from 'node-emoji';
import chalk from 'chalk';
import fs from 'fs-extra';

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
                        'chat config:set'
                    )} to generate a configuration file. ${emoji.get(
                        'pensive'
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
