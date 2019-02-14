"use strict";

var cov_2nmzwr8l86 = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/user/add.js";
  var hash = "d7ccfaca434bfd98254b82ab8856acda6f0cda11";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/user/add.js",
    statementMap: {
      "0": {
        start: {
          line: 30,
          column: 26
        },
        end: {
          line: 30,
          column: 45
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
          column: 28
        },
        end: {
          line: 38,
          column: 70
        }
      },
      "4": {
        start: {
          line: 39,
          column: 12
        },
        end: {
          line: 39,
          column: 69
        }
      },
      "5": {
        start: {
          line: 41,
          column: 12
        },
        end: {
          line: 46,
          column: 14
        }
      },
      "6": {
        start: {
          line: 47,
          column: 12
        },
        end: {
          line: 47,
          column: 25
        }
      },
      "7": {
        start: {
          line: 49,
          column: 12
        },
        end: {
          line: 49,
          column: 41
        }
      },
      "8": {
        start: {
          line: 54,
          column: 0
        },
        end: {
          line: 54,
          column: 53
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
exports.UserAdd = void 0;

var _command = require("@oclif/command");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class UserAdd extends _command.Command {
  async run() {
    cov_2nmzwr8l86.f[0]++;
    const {
      flags
    } = (cov_2nmzwr8l86.s[0]++, this.parse(UserAdd));
    cov_2nmzwr8l86.s[1]++;

    try {
      const client = (cov_2nmzwr8l86.s[2]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      const channel = (cov_2nmzwr8l86.s[3]++, await client.channel(flags.type, flags.id));
      cov_2nmzwr8l86.s[4]++;
      await channel.addModerators(flags.moderators.split(','));
      cov_2nmzwr8l86.s[5]++;
      this.log(`${flags.moderators} have been added as moderators to channel ${flags.type}:${flags.id}`, _nodeEmoji.default.get('rocket'));
      cov_2nmzwr8l86.s[6]++;
      this.exit(0);
    } catch (err) {
      cov_2nmzwr8l86.s[7]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.UserAdd = UserAdd;

_defineProperty(UserAdd, "flags", {
  id: _command.flags.string({
    char: 'i',
    description: _chalk.default.blue.bold('Channel name.'),
    required: true
  }),
  type: _command.flags.string({
    char: 't',
    description: _chalk.default.blue.bold('Channel type.'),
    required: true
  }),
  moderators: _command.flags.string({
    char: 'm',
    description: _chalk.default.blue.bold('Comma separated list of moderators to add.'),
    required: true
  })
});

cov_2nmzwr8l86.s[8]++;
UserAdd.description = 'Remove users from a channel.';