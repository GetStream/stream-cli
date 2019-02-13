import { expect, test } from '@oclif/test';

describe('channel', () => {
    test.stdout()
        .command(['channel:list'])
        .it('runs command:list', (ctx, done) => {
            expect(ctx.stdout).to.have.string('The channel');
            done();
        });
});
