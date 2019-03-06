const { expect, test } = require('@oclif/test');
const uuid = require('uuid/v4');

const channelId = uuid();

describe('channel', () => {
    test.stdout()
        .command([
            'chat:channel:create',
            `--channel=${channelId}`,
            '--type=messaging',
            '--name=CLI',
            '--image=https://images.unsplash.com/photo-1527427337751-fdca2f128ce5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=250&q=80',
            '--json',
        ])
        .exit(1)
        .it('runs chat:channel:create', ctx => {
            const data = JSON.parse(ctx.stdout);

            expect(data).to.be.an('object');
        });

    test.stdout()
        .command([
            'chat:channel:get',
            `--channel=${channelId}`,
            '--type=messaging',
            '--json',
        ])
        .exit(1)
        .it('runs chat:channel:get', ctx => {
            const data = JSON.parse(ctx.stdout);

            expect(data).to.be.an('object');
        });

    test.stdout()
        .command(['chat:channel:list', '--json'])
        .exit(1)
        .it('runs chat:channel:list', ctx => {
            const data = JSON.parse(ctx.stdout);

            expect(data).to.be.an('array');
        });

    test.stdout()
        .command(['chat:channel:query', `--channel=${channelId}`, '--json'])
        .exit(1)
        .it('runs chat:channel:query', ctx => {
            const data = JSON.parse(ctx.stdout);

            expect(data).to.be.an('object');
        });

    test.stdout()
        .command([
            'chat:channel:update',
            `--channel=${channelId}`,
            '--type=messaging',
            '--name=CLI',
            '--image=https://images.unsplash.com/photo-1527427337751-fdca2f128ce5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=250&q=80',
            '--reason=TEST',
            '--json',
        ])
        .exit(1)
        .it('runs chat:channel:update', ctx => {
            const data = JSON.parse(ctx.stdout);

            expect(data).to.be.an('object');
        });
});
