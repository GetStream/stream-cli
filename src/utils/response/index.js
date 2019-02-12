import chalk from 'chalk';
import emoji from 'node-emoji';
import moment from 'moment';

export function exit(msg, opts) {
    const ts = chalk.yellow.bold(
        moment().format('dddd, MMMM Do YYYY [at] h:mm:ss A')
    );

    if (opts.newline === true) {
        console.log(`${ts}:\n ${msg}`, opts.emoji ? emoji.get(opts.emoji) : '');
    } else {
        console.log(`${ts}:`, msg, emoji.get(opts.emoji || 'smile'));
    }

    process.exit(0);
}
