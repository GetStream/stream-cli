"use strict";

var cov_1ws4rif622 = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/user/remove.js";
  var hash = "5f123347e454bb862843bbbf902b29f38d2479fe";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/user/remove.js",
    statementMap: {
      "0": {
        start: {
          line: 32,
          column: 26
        },
        end: {
          line: 32,
          column: 48
        }
      },
      "1": {
        start: {
          line: 34,
          column: 8
        },
        end: {
          line: 52,
          column: 9
        }
      },
      "2": {
        start: {
          line: 35,
          column: 27
        },
        end: {
          line: 38,
          column: 13
        }
      },
      "3": {
        start: {
          line: 40,
          column: 28
        },
        end: {
          line: 40,
          column: 70
        }
      },
      "4": {
        start: {
          line: 41,
          column: 12
        },
        end: {
          line: 41,
          column: 72
        }
      },
      "5": {
        start: {
          line: 43,
          column: 12
        },
        end: {
          line: 48,
          column: 14
        }
      },
      "6": {
        start: {
          line: 49,
          column: 12
        },
        end: {
          line: 49,
          column: 25
        }
      },
      "7": {
        start: {
          line: 51,
          column: 12
        },
        end: {
          line: 51,
          column: 41
        }
      },
      "8": {
        start: {
          line: 56,
          column: 0
        },
        end: {
          line: 56,
          column: 56
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 31,
            column: 4
          },
          end: {
            line: 31,
            column: 5
          }
        },
        loc: {
          start: {
            line: 31,
            column: 16
          },
          end: {
            line: 53,
            column: 5
          }
        },
        line: 31
      }
    },
    branchMap: {},
    s: {
      "0": 0,
      "1": 0,
      "2": 0,
      "3": 0,
      "4": 0,
      "5": 0,
      "6": 0,
      "7": 0,
      "8": 0
    },
    f: {
      "0": 0
    },
    b: {},
    _coverageSchema: "43e27e138ebf9cfc5966b082cf9a028302ed4184"
  };
  var coverage = global[gcv] || (global[gcv] = {});

  if (coverage[path] && coverage[path].hash === hash) {
    return coverage[path];
  }

  coverageData.hash = hash;
  return coverage[path] = coverageData;
}();

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.UserRemove = void 0;

var _command = require("@oclif/command");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _uuid = _interopRequireDefault(require("uuid"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class UserRemove extends _command.Command {
  async run() {
    cov_1ws4rif622.f[0]++;
    const {
      flags
    } = (cov_1ws4rif622.s[0]++, this.parse(UserRemove));
    cov_1ws4rif622.s[1]++;

    try {
      const client = (cov_1ws4rif622.s[2]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      const channel = (cov_1ws4rif622.s[3]++, await client.channel(flags.type, flags.id));
      cov_1ws4rif622.s[4]++;
      await channel.demoteModerators(flags.moderators.split(','));
      cov_1ws4rif622.s[5]++;
      this.log(`${flags.moderators} have been removed as moderators from the ${flags.type} channel ${flags.id}`, _nodeEmoji.default.get('warning'));
      cov_1ws4rif622.s[6]++;
      this.exit(0);
    } catch (err) {
      cov_1ws4rif622.s[7]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.UserRemove = UserRemove;

_defineProperty(UserRemove, "flags", {
  id: _command.flags.string({
    char: 'i',
    description: _chalk.default.blue.bold('Channel name.'),
    default: (0, _uuid.default)(),
    required: true
  }),
  type: _command.flags.string({
    char: 't',
    description: _chalk.default.blue.bold('Channel type.'),
    required: true
  }),
  moderators: _command.flags.string({
    char: 'm',
    description: _chalk.default.blue.bold('Comma separated list of moderators to remove.'),
    required: true
  })
});

cov_1ws4rif622.s[8]++;
UserRemove.description = 'Remove users from a channel.';