"use strict";

var cov_1eaxcf0l45 = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/moderate/ban.js";
  var hash = "d4167f7f9ece14474759225b352f7a38f7528ad8";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/moderate/ban.js",
    statementMap: {
      "0": {
        start: {
          line: 30,
          column: 26
        },
        end: {
          line: 30,
          column: 49
        }
      },
      "1": {
        start: {
          line: 32,
          column: 8
        },
        end: {
          line: 50,
          column: 9
        }
      },
      "2": {
        start: {
          line: 33,
          column: 27
        },
        end: {
          line: 36,
          column: 13
        }
      },
      "3": {
        start: {
          line: 38,
          column: 12
        },
        end: {
          line: 41,
          column: 15
        }
      },
      "4": {
        start: {
          line: 43,
          column: 12
        },
        end: {
          line: 46,
          column: 14
        }
      },
      "5": {
        start: {
          line: 47,
          column: 12
        },
        end: {
          line: 47,
          column: 25
        }
      },
      "6": {
        start: {
          line: 49,
          column: 12
        },
        end: {
          line: 49,
          column: 41
        }
      },
      "7": {
        start: {
          line: 54,
          column: 0
        },
        end: {
          line: 54,
          column: 79
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 29,
            column: 4
          },
          end: {
            line: 29,
            column: 5
          }
        },
        loc: {
          start: {
            line: 29,
            column: 16
          },
          end: {
            line: 51,
            column: 5
          }
        },
        line: 29
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
      "7": 0
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
exports.ModerateBan = void 0;

var _command = require("@oclif/command");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class ModerateBan extends _command.Command {
  async run() {
    cov_1eaxcf0l45.f[0]++;
    const {
      flags
    } = (cov_1eaxcf0l45.s[0]++, this.parse(ModerateBan));
    cov_1eaxcf0l45.s[1]++;

    try {
      const client = (cov_1eaxcf0l45.s[2]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      cov_1eaxcf0l45.s[3]++;
      await client.banUser(flags.user, {
        timeout: Number(flags.timeout),
        reason: flags.reason
      });
      cov_1eaxcf0l45.s[4]++;
      this.log(`The user ${flags.user} has been banned!`, _nodeEmoji.default.get('banned'));
      cov_1eaxcf0l45.s[5]++;
      this.exit(0);
    } catch (err) {
      cov_1eaxcf0l45.s[6]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.ModerateBan = ModerateBan;

_defineProperty(ModerateBan, "flags", {
  user: _command.flags.string({
    char: 'u',
    description: _chalk.default.blue.bold('ID of user.'),
    exclusive: ['message'],
    required: true
  }),
  reason: _command.flags.string({
    char: 'r',
    description: _chalk.default.blue.bold('Reason for timeout.'),
    required: true
  }),
  timeout: _command.flags.string({
    char: 't',
    description: _chalk.default.blue.bold('Timeout in minutes.'),
    default: '60',
    required: true
  })
});

cov_1eaxcf0l45.s[7]++;
ModerateBan.description = 'Ban users indefinitely or by a per minute timeout.';