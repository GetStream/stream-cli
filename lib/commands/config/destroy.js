"use strict";

var cov_1bgc42zmjt = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/config/destroy.js";
  var hash = "59f40d50d862a1c26b2ea5fe6798fff73ef9624f";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/config/destroy.js",
    statementMap: {
      "0": {
        start: {
          line: 9,
          column: 8
        },
        end: {
          line: 22,
          column: 9
        }
      },
      "1": {
        start: {
          line: 10,
          column: 12
        },
        end: {
          line: 10,
          column: 77
        }
      },
      "2": {
        start: {
          line: 12,
          column: 12
        },
        end: {
          line: 17,
          column: 14
        }
      },
      "3": {
        start: {
          line: 19,
          column: 12
        },
        end: {
          line: 19,
          column: 25
        }
      },
      "4": {
        start: {
          line: 21,
          column: 12
        },
        end: {
          line: 21,
          column: 41
        }
      },
      "5": {
        start: {
          line: 26,
          column: 0
        },
        end: {
          line: 26,
          column: 45
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 8,
            column: 4
          },
          end: {
            line: 8,
            column: 5
          }
        },
        loc: {
          start: {
            line: 8,
            column: 16
          },
          end: {
            line: 23,
            column: 5
          }
        },
        line: 8
      }
    },
    branchMap: {},
    s: {
      "0": 0,
      "1": 0,
      "2": 0,
      "3": 0,
      "4": 0,
      "5": 0
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
exports.ConfigDestroy = void 0;

var _command = require("@oclif/command");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _fsExtra = _interopRequireDefault(require("fs-extra"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class ConfigDestroy extends _command.Command {
  async run() {
    cov_1bgc42zmjt.f[0]++;
    cov_1bgc42zmjt.s[0]++;

    try {
      cov_1bgc42zmjt.s[1]++;
      await _fsExtra.default.remove(_path.default.join(this.config.configDir, 'config.json'));
      cov_1bgc42zmjt.s[2]++;
      this.log(`Config destroyed. Run the command ${_chalk.default.blue.bold('config:set')} to generate a new config.`, _nodeEmoji.default.get('rocket'));
      cov_1bgc42zmjt.s[3]++;
      this.exit(0);
    } catch (err) {
      cov_1bgc42zmjt.s[4]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.ConfigDestroy = ConfigDestroy;
cov_1bgc42zmjt.s[5]++;
ConfigDestroy.description = 'Destroy config';