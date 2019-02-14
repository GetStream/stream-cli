"use strict";

var cov_18k8aj0gfc = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/channel/list.js";
  var hash = "37d7e65a1b00ab7255a97456223292338612a006";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/channel/list.js",
    statementMap: {
      "0": {
        start: {
          line: 10,
          column: 8
        },
        end: {
          line: 52,
          column: 9
        }
      },
      "1": {
        start: {
          line: 11,
          column: 27
        },
        end: {
          line: 14,
          column: 13
        }
      },
      "2": {
        start: {
          line: 16,
          column: 29
        },
        end: {
          line: 22,
          column: 13
        }
      },
      "3": {
        start: {
          line: 24,
          column: 12
        },
        end: {
          line: 49,
          column: 13
        }
      },
      "4": {
        start: {
          line: 25,
          column: 16
        },
        end: {
          line: 39,
          column: 18
        }
      },
      "5": {
        start: {
          line: 26,
          column: 20
        },
        end: {
          line: 38,
          column: 21
        }
      },
      "6": {
        start: {
          line: 41,
          column: 16
        },
        end: {
          line: 41,
          column: 29
        }
      },
      "7": {
        start: {
          line: 43,
          column: 16
        },
        end: {
          line: 46,
          column: 18
        }
      },
      "8": {
        start: {
          line: 48,
          column: 16
        },
        end: {
          line: 48,
          column: 29
        }
      },
      "9": {
        start: {
          line: 51,
          column: 12
        },
        end: {
          line: 51,
          column: 41
        }
      },
      "10": {
        start: {
          line: 56,
          column: 0
        },
        end: {
          line: 56,
          column: 46
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 9,
            column: 4
          },
          end: {
            line: 9,
            column: 5
          }
        },
        loc: {
          start: {
            line: 9,
            column: 16
          },
          end: {
            line: 53,
            column: 5
          }
        },
        line: 9
      },
      "1": {
        name: "(anonymous_1)",
        decl: {
          start: {
            line: 25,
            column: 29
          },
          end: {
            line: 25,
            column: 30
          }
        },
        loc: {
          start: {
            line: 26,
            column: 20
          },
          end: {
            line: 38,
            column: 21
          }
        },
        line: 26
      }
    },
    branchMap: {
      "0": {
        loc: {
          start: {
            line: 24,
            column: 12
          },
          end: {
            line: 49,
            column: 13
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 24,
            column: 12
          },
          end: {
            line: 49,
            column: 13
          }
        }, {
          start: {
            line: 24,
            column: 12
          },
          end: {
            line: 49,
            column: 13
          }
        }],
        line: 24
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
      "10": 0
    },
    f: {
      "0": 0,
      "1": 0
    },
    b: {
      "0": [0, 0]
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
exports.ChannelList = void 0;

var _command = require("@oclif/command");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ChannelList extends _command.Command {
  async run() {
    cov_18k8aj0gfc.f[0]++;
    cov_18k8aj0gfc.s[0]++;

    try {
      const client = (cov_18k8aj0gfc.s[1]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      const channels = (cov_18k8aj0gfc.s[2]++, await client.queryChannels({}, {
        last_message_at: -1
      }, {
        subscribe: false
      }));
      cov_18k8aj0gfc.s[3]++;

      if (channels.length) {
        cov_18k8aj0gfc.b[0][0]++;
        cov_18k8aj0gfc.s[4]++;
        channels.map(channel => {
          cov_18k8aj0gfc.f[1]++;
          cov_18k8aj0gfc.s[5]++;
          return this.log(_chalk.default.blue(`The channel ${_chalk.default.bold(channel.id)} of type ${_chalk.default.bold(channel.type)} with the CID of ${_chalk.default.bold(channel.cid)} has ${_chalk.default.bold(channel.data.members.length)} members.`));
        });
        cov_18k8aj0gfc.s[6]++;
        this.exit(0);
      } else {
        cov_18k8aj0gfc.b[0][1]++;
        cov_18k8aj0gfc.s[7]++;
        this.warn(`Your application does not have any channels.`, _nodeEmoji.default.get('pensive'));
        cov_18k8aj0gfc.s[8]++;
        this.exit(0);
      }
    } catch (err) {
      cov_18k8aj0gfc.s[9]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.ChannelList = ChannelList;
cov_18k8aj0gfc.s[10]++;
ChannelList.description = 'List all channels';