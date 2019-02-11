import chalk from 'chalk';
import emoji from 'node-emoji';

export function authError() {
    console.log(
        chalk.red(
            `Credentials not found. Run ${chalk.bold(
                'chat init'
            )} to generate a configuration file. ${emoji.get('pensive')}`
        )
    );

    process.exit(0);
}
