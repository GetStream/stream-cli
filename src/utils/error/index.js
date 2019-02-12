import chalk from 'chalk';
import emoji from 'node-emoji';
import moment from 'moment';

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

export function apiError(err) {
    let message = err.message || 'An unknown error has occurred.';

    const timestamp = chalk.yellow.bold(
        moment().format('dddd, MMMM Do YYYY [at] h:mm:ss A')
    );

    console.log(`${timestamp}: ${chalk.red(message)}.`, emoji.get('pensive'));

    process.exit(1);
}
