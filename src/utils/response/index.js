import chalk from 'chalk';
import emoji from 'node-emoji';

export function exit(_message, _emoji) {
    const timestamp = chalk.yellow.bold(
        moment().format('dddd, MMMM Do YYYY [at] h:mm:ss A')
    );

    console.log(`${timestamp}:`, _message, emoji.get(_emoji));

    process.exit(0);
}
