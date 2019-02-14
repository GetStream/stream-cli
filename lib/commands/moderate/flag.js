"use strict";

var cov_25bofqns1i = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/moderate/flag.js";
  var hash = "5676ed221a6e49a4b3e9de5092d16d9999ed42e2";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/moderate/flag.js",
    statementMap: {
      "0": {
        start: {
          line: 25,
          column: 26
        },
        end: {
          line: 25,
          column: 50
        }
      },
      "1": {
        start: {
          line: 27,
          column: 8
        },
        end: {
          line: 59,
          column: 9
        }
      },
      "2": {
        start: {
          line: 28,
          column: 27
        },
        end: {
          line: 31,
          column: 13
        }
      },
      "3": {
        start: {
          line: 33,
          column: 12
        },
        end: {
          line: 56,
          column: 13
        }
      },
      "4": {
        start: {
          line: 34,
          column: 16
        },
        end: {
          line: 34,
          column: 50
        }
      },
      "5": {
        start: {
          line: 36,
          column: 16
        },
        end: {
          line: 39,
          column: 18
        }
      },
      "6": {
        start: {
          line: 40,
          column: 16
        },
        end: {
          line: 40,
          column: 29
        }
      },
      "7": {
        start: {
          line: 41,
          column: 19
        },
        end: {
          line: 56,
          column: 13
        }
      },
      "8": {
        start: {
          line: 42,
          column: 16
        },
        end: {
          line: 42,
          column: 56
        }
      },
      "9": {
        start: {
          line: 44,
          column: 16
        },
        end: {
          line: 47,
          column: 18
        }
      },
      "10": {
        start: {
          line: 48,
          column: 16
        },
        end: {
          line: 48,
          column: 29
        }
      },
      "11": {
        start: {
          line: 50,
          column: 16
        },
        end: {
          line: 54,
          column: 18
        }
      },
      "12": {
        start: {
          line: 55,
          column: 16
        },
        end: {
          line: 55,
          column: 29
        }
      },
      "13": {
        start: {
          line: 58,
          column: 12
        },
        end: {
          line: 58,
          column: 41
        }
      },
      "14": {
        start: {
          line: 63,
          column: 0
        },
        end: {
          line: 63,
          column: 54
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 24,
            column: 4
          },
          end: {
            line: 24,
            column: 5
          }
        },
        loc: {
          start: {
            line: 24,
            column: 16
          },
          end: {
            line: 60,
            column: 5
          }
        },
        line: 24
      }
    },
    branchMap: {
      "0": {
        loc: {
          start: {
            line: 33,
            column: 12
          },
          end: {
            line: 56,
            column: 13
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 33,
            column: 12
          },
          end: {
            line: 56,
            column: 13
          }
        }, {
          start: {
            line: 33,
            column: 12
          },
          end: {
            line: 56,
            column: 13
          }
        }],
        line: 33
      },
      "1": {
        loc: {
          start: {
            line: 41,
            column: 19
          },
          end: {
            line: 56,
            column: 13
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 41,
            column: 19
          },
          end: {
            line: 56,
            column: 13
          }
        }, {
          start: {
            line: 41,
            column: 19
          },
          end: {
            line: 56,
            column: 13
          }
        }],
        line: 41
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
      "12": 0,
      "13": 0,
      "14": 0
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
exports.ModerateFlag = void 0;

var _command = require("@oclif/command");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class ModerateFlag extends _command.Command {
  async run() {
    cov_25bofqns1i.f[0]++;
    const {
      flags
    } = (cov_25bofqns1i.s[0]++, this.parse(ModerateFlag));
    cov_25bofqns1i.s[1]++;

    try {
      const client = (cov_25bofqns1i.s[2]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      cov_25bofqns1i.s[3]++;

      if (flags.user) {
        cov_25bofqns1i.b[0][0]++;
        cov_25bofqns1i.s[4]++;
        await client.flagUser(flags.user);
        cov_25bofqns1i.s[5]++;
        this.log(`The user ${flags.user} has been flagged!`, _nodeEmoji.default.get('bangbang'));
        cov_25bofqns1i.s[6]++;
        this.exit(0);
      } else {
        cov_25bofqns1i.b[0][1]++;
        cov_25bofqns1i.s[7]++;

        if (flags.message) {
          cov_25bofqns1i.b[1][0]++;
          cov_25bofqns1i.s[8]++;
          await client.flagMessage(flags.message);
          cov_25bofqns1i.s[9]++;
          this.log(`The message ${flags.user} has been flagged!`, _nodeEmoji.default.get('bangbang'));
          cov_25bofqns1i.s[10]++;
          this.exit(0);
        } else {
          cov_25bofqns1i.b[1][1]++;
          cov_25bofqns1i.s[11]++;
          this.warn(`Please pass a valid command. Use the command ${_chalk.default.blue.bold('moderate:flag --help')} for more information.`);
          cov_25bofqns1i.s[12]++;
          this.exit(0);
        }
      }
    } catch (err) {
      cov_25bofqns1i.s[13]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.ModerateFlag = ModerateFlag;

_defineProperty(ModerateFlag, "flags", {
  user: _command.flags.string({
    char: 'u',
    description: _chalk.default.blue.bold('ID of user.'),
    exclusive: ['message'],
    required: false
  }),
  message: _command.flags.string({
    char: 'm',
    description: _chalk.default.blue.bold('ID of message.'),
    exclusive: ['user'],
    required: false
  })
});

cov_25bofqns1i.s[14]++;
ModerateFlag.description = 'Flag users and messages.';