"use strict";

var cov_250gbfdhkw = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/channel/edit.js";
  var hash = "72a47cefe5946a7d2d8a5e7ebca49b9d89ade9da";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/channel/edit.js",
    statementMap: {
      "0": {
        start: {
          line: 50,
          column: 26
        },
        end: {
          line: 50,
          column: 49
        }
      },
      "1": {
        start: {
          line: 52,
          column: 8
        },
        end: {
          line: 85,
          column: 9
        }
      },
      "2": {
        start: {
          line: 53,
          column: 27
        },
        end: {
          line: 56,
          column: 13
        }
      },
      "3": {
        start: {
          line: 57,
          column: 28
        },
        end: {
          line: 57,
          column: 70
        }
      },
      "4": {
        start: {
          line: 59,
          column: 26
        },
        end: {
          line: 65,
          column: 13
        }
      },
      "5": {
        start: {
          line: 66,
          column: 12
        },
        end: {
          line: 66,
          column: 57
        }
      },
      "6": {
        start: {
          line: 66,
          column: 29
        },
        end: {
          line: 66,
          column: 57
        }
      },
      "7": {
        start: {
          line: 67,
          column: 12
        },
        end: {
          line: 67,
          column: 74
        }
      },
      "8": {
        start: {
          line: 67,
          column: 31
        },
        end: {
          line: 67,
          column: 74
        }
      },
      "9": {
        start: {
          line: 69,
          column: 12
        },
        end: {
          line: 72,
          column: 13
        }
      },
      "10": {
        start: {
          line: 70,
          column: 31
        },
        end: {
          line: 70,
          column: 53
        }
      },
      "11": {
        start: {
          line: 71,
          column: 16
        },
        end: {
          line: 71,
          column: 61
        }
      },
      "12": {
        start: {
          line: 74,
          column: 12
        },
        end: {
          line: 77,
          column: 15
        }
      },
      "13": {
        start: {
          line: 79,
          column: 12
        },
        end: {
          line: 82,
          column: 14
        }
      },
      "14": {
        start: {
          line: 84,
          column: 12
        },
        end: {
          line: 84,
          column: 41
        }
      },
      "15": {
        start: {
          line: 89,
          column: 0
        },
        end: {
          line: 89,
          column: 43
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 49,
            column: 4
          },
          end: {
            line: 49,
            column: 5
          }
        },
        loc: {
          start: {
            line: 49,
            column: 16
          },
          end: {
            line: 86,
            column: 5
          }
        },
        line: 49
      }
    },
    branchMap: {
      "0": {
        loc: {
          start: {
            line: 66,
            column: 12
          },
          end: {
            line: 66,
            column: 57
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 66,
            column: 12
          },
          end: {
            line: 66,
            column: 57
          }
        }, {
          start: {
            line: 66,
            column: 12
          },
          end: {
            line: 66,
            column: 57
          }
        }],
        line: 66
      },
      "1": {
        loc: {
          start: {
            line: 67,
            column: 12
          },
          end: {
            line: 67,
            column: 74
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 67,
            column: 12
          },
          end: {
            line: 67,
            column: 74
          }
        }, {
          start: {
            line: 67,
            column: 12
          },
          end: {
            line: 67,
            column: 74
          }
        }],
        line: 67
      },
      "2": {
        loc: {
          start: {
            line: 69,
            column: 12
          },
          end: {
            line: 72,
            column: 13
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 69,
            column: 12
          },
          end: {
            line: 72,
            column: 13
          }
        }, {
          start: {
            line: 69,
            column: 12
          },
          end: {
            line: 72,
            column: 13
          }
        }],
        line: 69
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
      "14": 0,
      "15": 0
    },
    f: {
      "0": 0
    },
    b: {
      "0": [0, 0],
      "1": [0, 0],
      "2": [0, 0]
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
exports.ChannelEdit = void 0;

var _command = require("@oclif/command");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _v = _interopRequireDefault(require("uuid/v4"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class ChannelEdit extends _command.Command {
  async run() {
    cov_250gbfdhkw.f[0]++;
    const {
      flags
    } = (cov_250gbfdhkw.s[0]++, this.parse(ChannelEdit));
    cov_250gbfdhkw.s[1]++;

    try {
      const client = (cov_250gbfdhkw.s[2]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      const channel = (cov_250gbfdhkw.s[3]++, await client.channel(flags.type, flags.id));
      let payload = (cov_250gbfdhkw.s[4]++, {
        name: flags.name,
        updated_by: {
          id: (0, _v.default)(),
          name: 'CLI'
        }
      });
      cov_250gbfdhkw.s[5]++;

      if (flags.image) {
        cov_250gbfdhkw.b[0][0]++;
        cov_250gbfdhkw.s[6]++;
        payload.image = flags.image;
      } else {
        cov_250gbfdhkw.b[0][1]++;
      }

      cov_250gbfdhkw.s[7]++;

      if (flags.members) {
        cov_250gbfdhkw.b[1][0]++;
        cov_250gbfdhkw.s[8]++;
        payload.members = flags.members.split(',');
      } else {
        cov_250gbfdhkw.b[1][1]++;
      }

      cov_250gbfdhkw.s[9]++;

      if (flags.data) {
        cov_250gbfdhkw.b[2][0]++;
        const parsed = (cov_250gbfdhkw.s[10]++, JSON.parse(flags.data));
        cov_250gbfdhkw.s[11]++;
        payload = Object.assign({}, payload, parsed);
      } else {
        cov_250gbfdhkw.b[2][1]++;
      }

      cov_250gbfdhkw.s[12]++;
      await channel.update(payload, {
        name: flags.name,
        text: flags.reason
      });
      cov_250gbfdhkw.s[13]++;
      this.log(`The channel ${flags.id} has been modified!`, _nodeEmoji.default.get('rocket'));
    } catch (err) {
      cov_250gbfdhkw.s[14]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.ChannelEdit = ChannelEdit;

_defineProperty(ChannelEdit, "flags", {
  id: _command.flags.string({
    char: 'i',
    description: _chalk.default.blue.bold('Channel ID.'),
    required: true
  }),
  type: _command.flags.string({
    char: 't',
    description: _chalk.default.blue.bold('Type of channel.'),
    options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
    required: true
  }),
  name: _command.flags.string({
    char: 'n',
    description: _chalk.default.blue.bold('Name of room.'),
    required: true
  }),
  url: _command.flags.string({
    char: 'u',
    description: _chalk.default.blue.bold('URL to channel image.'),
    required: false
  }),
  reason: _command.flags.string({
    char: 'r',
    description: _chalk.default.blue.bold('Reason for changing channel.'),
    required: true
  }),
  members: _command.flags.string({
    char: 'm',
    description: _chalk.default.blue.bold('Comma separated list of members.'),
    required: false
  }),
  data: _command.flags.string({
    char: 'd',
    description: _chalk.default.blue.bold('Additional data as a JSON payload.'),
    required: false
  })
});

cov_250gbfdhkw.s[15]++;
ChannelEdit.description = 'Edit a channel';