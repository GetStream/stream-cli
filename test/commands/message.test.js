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
            expect(JSON.parse(ctx.stdout)).to.be.an('array');
        });
});
