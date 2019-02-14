"use strict";

var cov_i783mv28z = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/channel/init.js";
  var hash = "7cbd0631434594c0d678c80d325e45b7d5c03b0c";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/channel/init.js",
    statementMap: {
      "0": {
        start: {
          line: 45,
          column: 26
        },
        end: {
          line: 45,
          column: 49
        }
      },
      "1": {
        start: {
          line: 47,
          column: 8
        },
        end: {
          line: 77,
          column: 9
        }
      },
      "2": {
        start: {
          line: 48,
          column: 27
        },
        end: {
          line: 51,
          column: 13
        }
      },
      "3": {
        start: {
          line: 53,
          column: 26
        },
        end: {
          line: 59,
          column: 13
        }
      },
      "4": {
        start: {
          line: 60,
          column: 12
        },
        end: {
          line: 60,
          column: 57
        }
      },
      "5": {
        start: {
          line: 60,
          column: 29
        },
        end: {
          line: 60,
          column: 57
        }
      },
      "6": {
        start: {
          line: 61,
          column: 12
        },
        end: {
          line: 61,
          column: 74
        }
      },
      "7": {
        start: {
          line: 61,
          column: 31
        },
        end: {
          line: 61,
          column: 74
        }
      },
      "8": {
        start: {
          line: 63,
          column: 12
        },
        end: {
          line: 66,
          column: 13
        }
      },
      "9": {
        start: {
          line: 64,
          column: 31
        },
        end: {
          line: 64,
          column: 53
        }
      },
      "10": {
        start: {
          line: 65,
          column: 16
        },
        end: {
          line: 65,
          column: 61
        }
      },
      "11": {
        start: {
          line: 68,
          column: 28
        },
        end: {
          line: 68,
          column: 79
        }
      },
      "12": {
        start: {
          line: 69,
          column: 12
        },
        end: {
          line: 69,
          column: 35
        }
      },
      "13": {
        start: {
          line: 71,
          column: 12
        },
        end: {
          line: 73,
          column: 15
        }
      },
      "14": {
        start: {
          line: 74,
          column: 12
        },
        end: {
          line: 74,
          column: 25
        }
      },
      "15": {
        start: {
          line: 76,
          column: 12
        },
        end: {
          line: 76,
          column: 41
        }
      },
      "16": {
        start: {
          line: 81,
          column: 0
        },
        end: {
          line: 81,
          column: 50
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 44,
            column: 4
          },
          end: {
            line: 44,
            column: 5
          }
        },
        loc: {
          start: {
            line: 44,
            column: 16
          },
          end: {
            line: 78,
            column: 5
          }
        },
        line: 44
      }
    },
    branchMap: {
      "0": {
        loc: {
          start: {
            line: 60,
            column: 12
          },
          end: {
            line: 60,
            column: 57
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 60,
            column: 12
          },
          end: {
            line: 60,
            column: 57
          }
        }, {
          start: {
            line: 60,
            column: 12
          },
          end: {
            line: 60,
            column: 57
          }
        }],
        line: 60
      },
      "1": {
        loc: {
          start: {
            line: 61,
            column: 12
          },
          end: {
            line: 61,
            column: 74
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 61,
            column: 12
          },
          end: {
            line: 61,
            column: 74
          }
        }, {
          start: {
            line: 61,
            column: 12
          },
          end: {
            line: 61,
            column: 74
          }
        }],
        line: 61
      },
      "2": {
        loc: {
          start: {
            line: 63,
            column: 12
          },
          end: {
            line: 66,
            column: 13
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 63,
            column: 12
          },
          end: {
            line: 66,
            column: 13
          }
        }, {
          start: {
            line: 63,
            column: 12
          },
          end: {
            line: 66,
            column: 13
          }
        }],
        line: 63
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
      "15": 0,
      "16": 0
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
exports.ChannelInit = void 0;

var _command = require("@oclif/command");

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _v = _interopRequireDefault(require("uuid/v4"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class ChannelInit extends _command.Command {
  async run() {
    cov_i783mv28z.f[0]++;
    const {
      flags
    } = (cov_i783mv28z.s[0]++, this.parse(ChannelInit));
    cov_i783mv28z.s[1]++;

    try {
      const client = (cov_i783mv28z.s[2]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      let payload = (cov_i783mv28z.s[3]++, {
        name: flags.name,
        created_by: {
          id: (0, _v.default)(),
          name: 'CLI'
        }
      });
      cov_i783mv28z.s[4]++;

      if (flags.image) {
        cov_i783mv28z.b[0][0]++;
        cov_i783mv28z.s[5]++;
        payload.image = flags.image;
      } else {
        cov_i783mv28z.b[0][1]++;
      }

      cov_i783mv28z.s[6]++;

      if (flags.members) {
        cov_i783mv28z.b[1][0]++;
        cov_i783mv28z.s[7]++;
        payload.members = flags.members.split(',');
      } else {
        cov_i783mv28z.b[1][1]++;
      }

      cov_i783mv28z.s[8]++;

      if (flags.data) {
        cov_i783mv28z.b[2][0]++;
        const parsed = (cov_i783mv28z.s[9]++, JSON.parse(flags.data));
        cov_i783mv28z.s[10]++;
        payload = Object.assign({}, payload, parsed);
      } else {
        cov_i783mv28z.b[2][1]++;
      }

      const channel = (cov_i783mv28z.s[11]++, await client.channel(flags.type, flags.id, payload));
      cov_i783mv28z.s[12]++;
      await channel.create();
      cov_i783mv28z.s[13]++;
      this.log(`The channel ${flags.name} has been initialized!`, {
        emoji: 'rocket'
      });
      cov_i783mv28z.s[14]++;
      this.exit(0);
    } catch (err) {
      cov_i783mv28z.s[15]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.ChannelInit = ChannelInit;

_defineProperty(ChannelInit, "flags", {
  id: _command.flags.string({
    char: 'i',
    description: _chalk.default.blue.bold('Channel ID.'),
    default: (0, _v.default)(),
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
  image: _command.flags.string({
    char: 'u',
    description: _chalk.default.blue.bold('URL to channel image.'),
    required: false
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

cov_i783mv28z.s[16]++;
ChannelInit.description = 'Initialize a channel.';