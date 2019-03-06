const { expect, test } = require('@oclif/test');
const uuid = require('uuid/v4');

let channelId;
let userId;

before('init', done => {
    test.stdout()
        .command(['chat:channel:list', '--json'])
        .exit(1)
        .it('runs chat:channel:list', ctx => {
            const data = JSON.parse(ctx.stdout);

            channelId = data[0].id;
        });

    test.stdout()
        .command([
            'chat:user:create',
            `--channel=${channelId}`,
            '--type=messaging',
            `--users=${uuid()}`,
            '--json',
        ])
        .exit(1)
        .it('runs chat:user:create', ctx => {
            const data = JSON.parse(ctx.stdout);

            userId = data[0].id;
        });

    done();
});

describe('create', () => {
    test.stdout()
        .command([
            'chat:message:create',
            `--channel=${channelId}`,
            '--type=messaging',
            `--user=${userId}`,
            '--name=Nick',
            '--message=buttercup',
            '--json',
        ])
        .exit(1)
        .it('runs chat:message:create', ctx => {
            const data = JSON.parse(ctx.stdout);

            expect(data).to.be.an('object');
        });

    test.stdout()
        .command([
            'chat:message:list',
            `--channel=${channelId}`,
            '--type=messaging',
            '--json',
        ])
        .exit(1)
        .it('runs chat:message:list', ctx => {
            const data = JSON.parse(ctx.stdout);

            expect(data).to.be.an('array');
        });
});
