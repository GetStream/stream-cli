const { expect, test } = require('@oclif/test');

describe('message', () => {
    test.stdout()
        .command([
            'chat:message:list',
            '--channel=godevs',
            '--type=messaging',
            '--json',
        ])
        .exit(1)
        .it('runs chat:message:list', ctx => {
            expect(ctx.stdout).to.be.a('string');
        });
});
