import { expect, test } from '@oclif/test';

describe('config', () => {
    test.stdout()
        .command(['config:get'])
        .it('returns an object with credentials', ctx => {
            expect(ctx.stdout).to.have.all.keys(['apiKey', 'apiSecret']);
        });
});
