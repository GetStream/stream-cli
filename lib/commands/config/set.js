"use strict";

var cov_1b9itfm7rp = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/config/set.js";
  var hash = "667c78613bb538e37c5109cb465aee0897da6005";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/config/set.js",
    statementMap: {
      "0": {
        start: {
          line: 10,
          column: 23
        },
        end: {
          line: 10,
          column: 70
        }
      },
      "1": {
        start: {
          line: 12,
          column: 8
        },
        end: {
          line: 59,
          column: 9
        }
      },
      "2": {
        start: {
          line: 13,
          column: 27
        },
        end: {
          line: 13,
          column: 54
        }
      },
      "3": {
        start: {
          line: 15,
          column: 12
        },
        end: {
          line: 29,
          column: 13
        }
      },
      "4": {
        start: {
          line: 16,
          column: 31
        },
        end: {
          line: 24,
          column: 18
        }
      },
      "5": {
        start: {
          line: 26,
          column: 16
        },
        end: {
          line: 28,
          column: 17
        }
      },
      "6": {
        start: {
          line: 27,
          column: 20
        },
        end: {
          line: 27,
          column: 33
        }
      },
      "7": {
        start: {
          line: 31,
          column: 25
        },
        end: {
          line: 46,
          column: 14
        }
      },
      "8": {
        start: {
          line: 48,
          column: 12
        },
        end: {
          line: 51,
          column: 15
        }
      },
      "9": {
        start: {
          line: 53,
          column: 12
        },
        end: {
          line: 56,
          column: 14
        }
      },
      "10": {
        start: {
          line: 58,
          column: 12
        },
        end: {
          line: 58,
          column: 41
        }
      },
      "11": {
        start: {
          line: 63,
          column: 0
        },
        end: {
          line: 63,
          column: 65
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
            line: 60,
            column: 5
          }
        },
        line: 9
      }
    },
    branchMap: {
      "0": {
        loc: {
          start: {
            line: 15,
            column: 12
          },
          end: {
            line: 29,
            column: 13
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 15,
            column: 12
          },
          end: {
            line: 29,
            column: 13
          }
        }, {
          start: {
            line: 15,
            column: 12
          },
          end: {
            line: 29,
            column: 13
          }
        }],
        line: 15
      },
      "1": {
        loc: {
          start: {
            line: 26,
            column: 16
          },
          end: {
            line: 28,
            column: 17
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 26,
            column: 16
          },
          end: {
            line: 28,
            column: 17
          }
        }, {
          start: {
            line: 26,
            column: 16
          },
          end: {
            line: 28,
            column: 17
          }
        }],
        line: 26
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
      "11": 0
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
exports.ConfigSet = void 0;

var _command = require("@oclif/command");

var _enquirer = require("enquirer");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _fsExtra = _interopRequireDefault(require("fs-extra"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ConfigSet extends _command.Command {
  async run() {
    cov_1b9itfm7rp.f[0]++;
    const config = (cov_1b9itfm7rp.s[0]++, _path.default.join(this.config.configDir, 'config.json'));
    cov_1b9itfm7rp.s[1]++;

    try {
      const exists = (cov_1b9itfm7rp.s[2]++, await _fsExtra.default.pathExists(config));
      cov_1b9itfm7rp.s[3]++;

      if (exists) {
        cov_1b9itfm7rp.b[0][0]++;
        const answer = (cov_1b9itfm7rp.s[4]++, await (0, _enquirer.prompt)({
          type: 'confirm',
          name: 'continue',
          message: _chalk.default.yellow.bold(`This command will delete your current configuration. Do you want to continue? ${_nodeEmoji.default.get('warning')} `)
        }));
        cov_1b9itfm7rp.s[5]++;

        if (!answer.continue) {
          cov_1b9itfm7rp.b[1][0]++;
          cov_1b9itfm7rp.s[6]++;
          this.exit(0);
        } else {
          cov_1b9itfm7rp.b[1][1]++;
        }
      } else {
        cov_1b9itfm7rp.b[0][1]++;
      }

      const data = (cov_1b9itfm7rp.s[7]++, await (0, _enquirer.prompt)([{
        type: 'input',
        name: 'apiKey',
        message: _chalk.default.green(`What's your API key? ${_nodeEmoji.default.get('lock')}`)
      }, {
        type: 'input',
        name: 'apiSecret',
        message: _chalk.default.green(`What's your API secret? ${_nodeEmoji.default.get('lock')}`)
      }]));
      cov_1b9itfm7rp.s[8]++;
      await _fsExtra.default.writeJson(config, {
        apiKey: data.apiKey,
        apiSecret: data.apiSecret
      });
      cov_1b9itfm7rp.s[9]++;
      this.log(_chalk.default.green(`Your config has been generated!`), _nodeEmoji.default.get('rocket'));
    } catch (err) {
      cov_1b9itfm7rp.s[10]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.ConfigSet = ConfigSet;
cov_1b9itfm7rp.s[11]++;
ConfigSet.description = 'Stores your Stream API key and secret.';