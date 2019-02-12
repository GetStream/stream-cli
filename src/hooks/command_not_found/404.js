import chalk from 'chalk';
import emoji from 'node-emoji';

module.exports = async function(opts) {
    process.stdout.write(
        chalk.red(
            `Command not found. Please run ${chalk.bold(
                'chat help'
            )} for more information on commands. ${emoji.get('pensive')}\n`
        )
    );

    process.exit(1);
};
