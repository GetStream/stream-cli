"use strict";

var cov_1padvcu0s = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/channel/query.js";
  var hash = "a9d5054f334a19ebb3ce10caed4ce2b4089bb934";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/channel/query.js",
    statementMap: {
      "0": {
        start: {
          line: 35,
          column: 26
        },
        end: {
          line: 35,
          column: 50
        }
      },
      "1": {
        start: {
          line: 37,
          column: 8
        },
        end: {
          line: 55,
          column: 9
        }
      },
      "2": {
        start: {
          line: 38,
          column: 27
        },
        end: {
          line: 41,
          column: 13
        }
      },
      "3": {
        start: {
          line: 43,
          column: 27
        },
        end: {
          line: 43,
          column: 73
        }
      },
      "4": {
        start: {
          line: 44,
          column: 25
        },
        end: {
          line: 44,
          column: 65
        }
      },
      "5": {
        start: {
          line: 46,
          column: 29
        },
        end: {
          line: 48,
          column: 14
        }
      },
      "6": {
        start: {
          line: 50,
          column: 12
        },
        end: {
          line: 50,
          column: 39
        }
      },
      "7": {
        start: {
          line: 52,
          column: 12
        },
        end: {
          line: 52,
          column: 25
        }
      },
      "8": {
        start: {
          line: 54,
          column: 12
        },
        end: {
          line: 54,
          column: 41
        }
      },
      "9": {
        start: {
          line: 59,
          column: 0
        },
        end: {
          line: 59,
          column: 45
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 34,
            column: 4
          },
          end: {
            line: 34,
            column: 5
          }
        },
        loc: {
          start: {
            line: 34,
            column: 16
          },
          end: {
            line: 56,
            column: 5
          }
        },
        line: 34
      }
    },
    branchMap: {
      "0": {
        loc: {
          start: {
            line: 43,
            column: 27
          },
          end: {
            line: 43,
            column: 73
          }
        },
        type: "cond-expr",
        locations: [{
          start: {
            line: 43,
            column: 43
          },
          end: {
            line: 43,
            column: 68
          }
        }, {
          start: {
            line: 43,
            column: 71
          },
          end: {
            line: 43,
            column: 73
          }
        }],
        line: 43
      },
      "1": {
        loc: {
          start: {
            line: 44,
            column: 25
          },
          end: {
            line: 44,
            column: 65
          }
        },
        type: "cond-expr",
        locations: [{
          start: {
            line: 44,
            column: 38
          },
          end: {
            line: 44,
            column: 60
          }
        }, {
          start: {
            line: 44,
            column: 63
          },
          end: {
            line: 44,
            column: 65
          }
        }],
        line: 44
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
      "9": 0
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
exports.ChannelQuery = void 0;

var _command = require("@oclif/command");

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _v = _interopRequireDefault(require("uuid/v4"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class ChannelQuery extends _command.Command {
  async run() {
    cov_1padvcu0s.f[0]++;
    const {
      flags
    } = (cov_1padvcu0s.s[0]++, this.parse(ChannelQuery));
    cov_1padvcu0s.s[1]++;

    try {
      const client = (cov_1padvcu0s.s[2]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      const filter = (cov_1padvcu0s.s[3]++, flags.filters ? (cov_1padvcu0s.b[0][0]++, JSON.parse(flags.filters)) : (cov_1padvcu0s.b[0][1]++, {}));
      const sort = (cov_1padvcu0s.s[4]++, flags.sort ? (cov_1padvcu0s.b[1][0]++, JSON.parse(flags.sort)) : (cov_1padvcu0s.b[1][1]++, {}));
      const channels = (cov_1padvcu0s.s[5]++, await client.queryChannels(filter, sort, {
        subscribe: false
      }));
      cov_1padvcu0s.s[6]++;
      this.log(channels[0].data);
      cov_1padvcu0s.s[7]++;
      this.exit(0);
    } catch (err) {
      cov_1padvcu0s.s[8]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.ChannelQuery = ChannelQuery;

_defineProperty(ChannelQuery, "flags", {
  id: _command.flags.string({
    char: 'i',
    description: _chalk.default.blue.bold('Channel ID.'),
    default: (0, _v.default)(),
    required: false
  }),
  type: _command.flags.string({
    char: 't',
    description: _chalk.default.blue.bold('Type of channel.'),
    options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
    required: false
  }),
  filter: _command.flags.string({
    char: 'f',
    description: _chalk.default.blue.bold('Filters to apply.'),
    required: false
  }),
  sort: _command.flags.string({
    char: 's',
    description: _chalk.default.blue.bold('Sort to apply.'),
    required: false
  })
});

cov_1padvcu0s.s[9]++;
ChannelQuery.description = 'Query a channel';