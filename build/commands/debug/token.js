"use strict";

var _command = require("@oclif/command");

var _cliTable = _interopRequireDefault(require("cli-table"));

var _enquirer = require("enquirer");

var _chalk = _interopRequireDefault(require("chalk"));

var _jsonwebtoken = _interopRequireDefault(require("jsonwebtoken"));

var _config = require("../../utils/config");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class DebugToken extends _command.Command {
  async run() {
    const {
      flags
    } = this.parse(DebugToken);

    try {
      if (!flags.token) {
        const res = await (0, _enquirer.prompt)([{
          type: 'input',
          name: 'jwt',
          message: `What is the Stream token you would like to debug?`,
          required: true
        }]);
        flags.jwt = res.jwt;
      }

      const {
        apiSecret
      } = await (0, _config.credentials)(this);
      const decoded = await _jsonwebtoken.default.verify(flags.jwt, apiSecret, {
        complete: true
      });

      if (flags.json) {
        this.log(JSON.stringify(decoded));
        this.exit(0);
      }

      const table = new _cliTable.default();
      table.push({
        [`${_chalk.default.green.bold('Header Type')}`]: decoded.header.typ
      }, {
        [`${_chalk.default.green.bold('Header Algorithm')}`]: decoded.header.alg
      }, {
        [`${_chalk.default.green.bold('Signature')}`]: decoded.signature
      }, {
        [`${_chalk.default.green.bold('User ID')}`]: decoded.payload.user_id
      });
      this.log(table.toString());
      this.exit(0);
    } catch (error) {
      this.error('Malformed JWT token or Stream API secret.', {
        exit: 1
      });
    }
  }

}

DebugToken.flags = {
  token: _command.flags.string({
    char: 't',
    description: 'The Stream token you are trying to debug.',
    required: false
  }),
  json: _command.flags.boolean({
    char: 'j',
    description: 'Output results in JSON. When not specified, returns output in a human friendly format.',
    required: false
  })
};
DebugToken.description = 'Debugs a JWT token provided by Stream.';
module.exports.DebugToken = DebugToken;