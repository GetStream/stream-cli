"use strict";

var cov_sebwynz2l = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/config/get.js";
  var hash = "41866192b0ccf868f00ce97881f08b47d9a9dfe4";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/config/get.js",
    statementMap: {
      "0": {
        start: {
          line: 11,
          column: 23
        },
        end: {
          line: 11,
          column: 70
        }
      },
      "1": {
        start: {
          line: 12,
          column: 38
        },
        end: {
          line: 12,
          column: 69
        }
      },
      "2": {
        start: {
          line: 14,
          column: 8
        },
        end: {
          line: 35,
          column: 9
        }
      },
      "3": {
        start: {
          line: 15,
          column: 26
        },
        end: {
          line: 18,
          column: 14
        }
      },
      "4": {
        start: {
          line: 20,
          column: 12
        },
        end: {
          line: 20,
          column: 44
        }
      },
      "5": {
        start: {
          line: 22,
          column: 12
        },
        end: {
          line: 22,
          column: 39
        }
      },
      "6": {
        start: {
          line: 23,
          column: 12
        },
        end: {
          line: 23,
          column: 25
        }
      },
      "7": {
        start: {
          line: 25,
          column: 12
        },
        end: {
          line: 34,
          column: 14
        }
      },
      "8": {
        start: {
          line: 39,
          column: 0
        },
        end: {
          line: 39,
          column: 68
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 10,
            column: 4
          },
          end: {
            line: 10,
            column: 5
          }
        },
        loc: {
          start: {
            line: 10,
            column: 16
          },
          end: {
            line: 36,
            column: 5
          }
        },
        line: 10
      }
    },
    branchMap: {
      "0": {
        loc: {
          start: {
            line: 14,
            column: 8
          },
          end: {
            line: 35,
            column: 9
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 14,
            column: 8
          },
          end: {
            line: 35,
            column: 9
          }
        }, {
          start: {
            line: 14,
            column: 8
          },
          end: {
            line: 35,
            column: 9
          }
        }],
        line: 14
      },
      "1": {
        loc: {
          start: {
            line: 14,
            column: 12
          },
          end: {
            line: 14,
            column: 31
          }
        },
        type: "binary-expr",
        locations: [{
          start: {
            line: 14,
            column: 12
          },
          end: {
            line: 14,
            column: 18
          }
        }, {
          start: {
            line: 14,
            column: 22
          },
          end: {
            line: 14,
            column: 31
          }
        }],
        line: 14
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
      "8": 0
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
exports.ConfigGet = void 0;

var _command = require("@oclif/command");

var _cliTable = _interopRequireDefault(require("cli-table"));

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _config = require("../../utils/config");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ConfigGet extends _command.Command {
  async run() {
    cov_sebwynz2l.f[0]++;
    const config = (cov_sebwynz2l.s[0]++, _path.default.join(this.config.configDir, 'config.json'));
    const {
      apiKey,
      apiSecret
    } = (cov_sebwynz2l.s[1]++, await (0, _config.credentials)(config, this));
    cov_sebwynz2l.s[2]++;

    if ((cov_sebwynz2l.b[1][0]++, apiKey) && (cov_sebwynz2l.b[1][1]++, apiSecret)) {
      cov_sebwynz2l.b[0][0]++;
      const table = (cov_sebwynz2l.s[3]++, new _cliTable.default({
        head: ['API Key', 'API Secret'],
        colWidths: [25, 75]
      }));
      cov_sebwynz2l.s[4]++;
      table.push([apiKey, apiSecret]);
      cov_sebwynz2l.s[5]++;
      this.log(table.toString());
      cov_sebwynz2l.s[6]++;
      this.exit(0);
    } else {
      cov_sebwynz2l.b[0][1]++;
      cov_sebwynz2l.s[7]++;
      this.error(_chalk.default.red(`Credentials not found. Run ${_chalk.default.bold('chat init')} to generate a configuration file. ${_nodeEmoji.default.get('pensive')}`), {
        exit: 1
      });
    }
  }

}

exports.ConfigGet = ConfigGet;
cov_sebwynz2l.s[8]++;
ConfigGet.description = 'Retrieves API config credentials for CLI.';