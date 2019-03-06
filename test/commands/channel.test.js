const { expect, test } = require('@oclif/test');
const uuid = require('uuid/v4');

describe('channel', () => {
    test.stdout()
        .command([
            'chat:channel:create',
            `--channel=${uuid()}`,
            '--type=messaging',
            '--name=CLI',
            '--image=https://images.unsplash.com/photo-1527427337751-fdca2f128ce5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=250&q=80',
            '--json',
        ])
        .exit(1)
        .it('runs chat:channel:create', (ctx, done) => {
            const data = JSON.parse(ctx.stdout);

            expect(data).to.be.an('object');
            expect(data).to.have.property('channel');

            done();
        });
});
