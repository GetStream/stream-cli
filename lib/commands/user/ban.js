"use strict";

var cov_1pbtoqqgya = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/user/ban.js";
  var hash = "d44e217217c5eff3fa1a1c24743a8bc9229212cb";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/user/ban.js",
    statementMap: {
      "0": {
        start: {
          line: 42,
          column: 26
        },
        end: {
          line: 42,
          column: 45
        }
      },
      "1": {
        start: {
          line: 44,
          column: 8
        },
        end: {
          line: 65,
          column: 9
        }
      },
      "2": {
        start: {
          line: 45,
          column: 27
        },
        end: {
          line: 48,
          column: 13
        }
      },
      "3": {
        start: {
          line: 50,
          column: 28
        },
        end: {
          line: 50,
          column: 30
        }
      },
      "4": {
        start: {
          line: 51,
          column: 12
        },
        end: {
          line: 51,
          column: 63
        }
      },
      "5": {
        start: {
          line: 51,
          column: 31
        },
        end: {
          line: 51,
          column: 63
        }
      },
      "6": {
        start: {
          line: 52,
          column: 12
        },
        end: {
          line: 52,
          column: 60
        }
      },
      "7": {
        start: {
          line: 52,
          column: 30
        },
        end: {
          line: 52,
          column: 60
        }
      },
      "8": {
        start: {
          line: 54,
          column: 12
        },
        end: {
          line: 54,
          column: 54
        }
      },
      "9": {
        start: {
          line: 56,
          column: 12
        },
        end: {
          line: 61,
          column: 14
        }
      },
      "10": {
        start: {
          line: 62,
          column: 12
        },
        end: {
          line: 62,
          column: 25
        }
      },
      "11": {
        start: {
          line: 64,
          column: 12
        },
        end: {
          line: 64,
          column: 41
        }
      },
      "12": {
        start: {
          line: 69,
          column: 0
        },
        end: {
          line: 69,
          column: 74
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 41,
            column: 4
          },
          end: {
            line: 41,
            column: 5
          }
        },
        loc: {
          start: {
            line: 41,
            column: 16
          },
          end: {
            line: 66,
            column: 5
          }
        },
        line: 41
      }
    },
    branchMap: {
      "0": {
        loc: {
          start: {
            line: 51,
            column: 12
          },
          end: {
            line: 51,
            column: 63
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 51,
            column: 12
          },
          end: {
            line: 51,
            column: 63
          }
        }, {
          start: {
            line: 51,
            column: 12
          },
          end: {
            line: 51,
            column: 63
          }
        }],
        line: 51
      },
      "1": {
        loc: {
          start: {
            line: 52,
            column: 12
          },
          end: {
            line: 52,
            column: 60
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 52,
            column: 12
          },
          end: {
            line: 52,
            column: 60
          }
        }, {
          start: {
            line: 52,
            column: 12
          },
          end: {
            line: 52,
            column: 60
          }
        }],
        line: 52
      }
    },
    s: {
      "0": 0,
      "1": 0,
      "2": 0,
      "3": 0,
      "4": 0,
      "5": 0,
      "6": 0,
      "7": 0,
      "8": 0,
      "9": 0,
      "10": 0,
      "11": 0,
      "12": 0
    },
    f: {
      "0": 0
    },
    b: {
      "0": [0, 0],
      "1": [0, 0]
    },
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
exports.UserBan = void 0;

var _command = require("@oclif/command");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _uuid = _interopRequireDefault(require("uuid"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class UserBan extends _command.Command {
  async run() {
    cov_1pbtoqqgya.f[0]++;
    const {
      flags
    } = (cov_1pbtoqqgya.s[0]++, this.parse(UserBan));
    cov_1pbtoqqgya.s[1]++;

    try {
      const client = (cov_1pbtoqqgya.s[2]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      const payload = (cov_1pbtoqqgya.s[3]++, {});
      cov_1pbtoqqgya.s[4]++;

      if (flags.timeout) {
        cov_1pbtoqqgya.b[0][0]++;
        cov_1pbtoqqgya.s[5]++;
        payload.timeout = flags.timeout;
      } else {
        cov_1pbtoqqgya.b[0][1]++;
      }

      cov_1pbtoqqgya.s[6]++;

      if (flags.reason) {
        cov_1pbtoqqgya.b[1][0]++;
        cov_1pbtoqqgya.s[7]++;
        payload.reason = flags.reason;
      } else {
        cov_1pbtoqqgya.b[1][1]++;
      }

      cov_1pbtoqqgya.s[8]++;
      await client.banUser(flags.user, payload);
      cov_1pbtoqqgya.s[9]++;
      this.log(`${flags.user} has been added banned from channel ${flags.type}:${flags.id}`, _nodeEmoji.default.get('banned'));
      cov_1pbtoqqgya.s[10]++;
      this.exit(0);
    } catch (err) {
      cov_1pbtoqqgya.s[11]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.UserBan = UserBan;

_defineProperty(UserBan, "flags", {
  id: _command.flags.string({
    char: 'i',
    description: _chalk.default.blue.bold('Channel ID.'),
    default: (0, _uuid.default)(),
    required: true
  }),
  type: _command.flags.string({
    char: 't',
    description: _chalk.default.blue.bold('Type of channel.'),
    options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
    required: true
  }),
  user: _command.flags.string({
    char: 'u',
    description: _chalk.default.blue.bold('User ID.'),
    required: true
  }),
  reason: _command.flags.string({
    char: 'r',
    description: _chalk.default.blue.bold('Reason to place ban.'),
    required: false
  }),
  timeout: _command.flags.string({
    char: 't',
    description: _chalk.default.blue.bold('Duration in minutes.'),
    default: '60',
    required: false
  })
});

cov_1pbtoqqgya.s[12]++;
UserBan.description = 'Ban users indefinitely or by a per-minute period.';