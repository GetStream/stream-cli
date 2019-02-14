"use strict";

var cov_76s8w6qau = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/message/remove.js";
  var hash = "0f9cab5246ddba6df83463361bc5c46cee1e0872";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/message/remove.js",
    statementMap: {
      "0": {
        start: {
          line: 20,
          column: 26
        },
        end: {
          line: 20,
          column: 51
        }
      },
      "1": {
        start: {
          line: 22,
          column: 8
        },
        end: {
          line: 37,
          column: 9
        }
      },
      "2": {
        start: {
          line: 23,
          column: 27
        },
        end: {
          line: 26,
          column: 13
        }
      },
      "3": {
        start: {
          line: 28,
          column: 12
        },
        end: {
          line: 28,
          column: 49
        }
      },
      "4": {
        start: {
          line: 30,
          column: 12
        },
        end: {
          line: 33,
          column: 14
        }
      },
      "5": {
        start: {
          line: 34,
          column: 12
        },
        end: {
          line: 34,
          column: 25
        }
      },
      "6": {
        start: {
          line: 36,
          column: 12
        },
        end: {
          line: 36,
          column: 41
        }
      },
      "7": {
        start: {
          line: 41,
          column: 0
        },
        end: {
          line: 41,
          column: 58
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 19,
            column: 4
          },
          end: {
            line: 19,
            column: 5
          }
        },
        loc: {
          start: {
            line: 19,
            column: 16
          },
          end: {
            line: 38,
            column: 5
          }
        },
        line: 19
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
exports.MessageRemove = void 0;

var _command = require("@oclif/command");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _uuid = _interopRequireDefault(require("uuid"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class MessageRemove extends _command.Command {
  async run() {
    cov_76s8w6qau.f[0]++;
    const {
      flags
    } = (cov_76s8w6qau.s[0]++, this.parse(MessageRemove));
    cov_76s8w6qau.s[1]++;

    try {
      const client = (cov_76s8w6qau.s[2]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      cov_76s8w6qau.s[3]++;
      await client.deleteMessage(flags.id);
      cov_76s8w6qau.s[4]++;
      this.log(`The message ${flags.id} has been removed!`, _nodeEmoji.default.get('wastebasket'));
      cov_76s8w6qau.s[5]++;
      this.exit(0);
    } catch (err) {
      cov_76s8w6qau.s[6]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.MessageRemove = MessageRemove;

_defineProperty(MessageRemove, "flags", {
  id: _command.flags.string({
    char: 'i',
    description: _chalk.default.blue.bold('Channel ID.'),
    default: (0, _uuid.default)(),
    required: false
  })
});

cov_76s8w6qau.s[7]++;
MessageRemove.description = 'Send messages to a channel.';