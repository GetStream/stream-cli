import { expect, test } from '@oclif/test';

describe.skip('channel', () => {
    test.stdout()
        .command(['channel:list'])
        .it('returns a table with api credentials', ctx => {
            expect(ctx.stdout).to.have.all.keys(['apiKey', 'apiSecret']);
        });
});
