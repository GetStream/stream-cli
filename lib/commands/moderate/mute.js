"use strict";

var cov_2pazujm15j = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/moderate/mute.js";
  var hash = "89738792c3d1a24f629a7ef0806ebf0f42e9aeb5";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/moderate/mute.js",
    statementMap: {
      "0": {
        start: {
          line: 18,
          column: 26
        },
        end: {
          line: 18,
          column: 50
        }
      },
      "1": {
        start: {
          line: 20,
          column: 8
        },
        end: {
          line: 35,
          column: 9
        }
      },
      "2": {
        start: {
          line: 21,
          column: 27
        },
        end: {
          line: 24,
          column: 13
        }
      },
      "3": {
        start: {
          line: 26,
          column: 12
        },
        end: {
          line: 26,
          column: 46
        }
      },
      "4": {
        start: {
          line: 28,
          column: 12
        },
        end: {
          line: 31,
          column: 14
        }
      },
      "5": {
        start: {
          line: 32,
          column: 12
        },
        end: {
          line: 32,
          column: 25
        }
      },
      "6": {
        start: {
          line: 34,
          column: 12
        },
        end: {
          line: 34,
          column: 41
        }
      },
      "7": {
        start: {
          line: 39,
          column: 0
        },
        end: {
          line: 39,
          column: 58
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 17,
            column: 4
          },
          end: {
            line: 17,
            column: 5
          }
        },
        loc: {
          start: {
            line: 17,
            column: 16
          },
          end: {
            line: 36,
            column: 5
          }
        },
        line: 17
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
exports.ModerateMute = void 0;

var _command = require("@oclif/command");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class ModerateMute extends _command.Command {
  async run() {
    cov_2pazujm15j.f[0]++;
    const {
      flags
    } = (cov_2pazujm15j.s[0]++, this.parse(ModerateMute));
    cov_2pazujm15j.s[1]++;

    try {
      const client = (cov_2pazujm15j.s[2]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      cov_2pazujm15j.s[3]++;
      await client.muteUser(flags.user);
      cov_2pazujm15j.s[4]++;
      this.log(`The message ${flags.user} has been flagged!`, _nodeEmoji.default.get('two_flags'));
      cov_2pazujm15j.s[5]++;
      this.exit(0);
    } catch (err) {
      cov_2pazujm15j.s[6]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.ModerateMute = ModerateMute;

_defineProperty(ModerateMute, "flags", {
  user: _command.flags.string({
    char: 'u',
    description: _chalk.default.blue.bold('User ID.'),
    required: true
  })
});

cov_2pazujm15j.s[7]++;
ModerateMute.description = 'Mute users who are annoying.';