"use strict";

var cov_1iua2gc7za = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/utils/config/index.js";
  var hash = "fd47654abfcec9d94f98397d9492e75d9e2cd0ca";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/utils/config/index.js",
    statementMap: {
      "0": {
        start: {
          line: 6,
          column: 4
        },
        end: {
          line: 34,
          column: 5
        }
      },
      "1": {
        start: {
          line: 7,
          column: 8
        },
        end: {
          line: 12,
          column: 9
        }
      },
      "2": {
        start: {
          line: 8,
          column: 12
        },
        end: {
          line: 11,
          column: 15
        }
      },
      "3": {
        start: {
          line: 14,
          column: 38
        },
        end: {
          line: 14,
          column: 63
        }
      },
      "4": {
        start: {
          line: 16,
          column: 8
        },
        end: {
          line: 28,
          column: 9
        }
      },
      "5": {
        start: {
          line: 17,
          column: 12
        },
        end: {
          line: 25,
          column: 14
        }
      },
      "6": {
        start: {
          line: 27,
          column: 12
        },
        end: {
          line: 27,
          column: 26
        }
      },
      "7": {
        start: {
          line: 30,
          column: 8
        },
        end: {
          line: 30,
          column: 37
        }
      },
      "8": {
        start: {
          line: 32,
          column: 8
        },
        end: {
          line: 32,
          column: 25
        }
      },
      "9": {
        start: {
          line: 33,
          column: 8
        },
        end: {
          line: 33,
          column: 22
        }
      }
    },
    fnMap: {
      "0": {
        name: "credentials",
        decl: {
          start: {
            line: 5,
            column: 22
          },
          end: {
            line: 5,
            column: 33
          }
        },
        loc: {
          start: {
            line: 5,
            column: 49
          },
          end: {
            line: 35,
            column: 1
          }
        },
        line: 5
      }
    },
    branchMap: {
      "0": {
        loc: {
          start: {
            line: 7,
            column: 8
          },
          end: {
            line: 12,
            column: 9
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 7,
            column: 8
          },
          end: {
            line: 12,
            column: 9
          }
        }, {
          start: {
            line: 7,
            column: 8
          },
          end: {
            line: 12,
            column: 9
          }
        }],
        line: 7
      },
      "1": {
        loc: {
          start: {
            line: 16,
            column: 8
          },
          end: {
            line: 28,
            column: 9
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 16,
            column: 8
          },
          end: {
            line: 28,
            column: 9
          }
        }, {
          start: {
            line: 16,
            column: 8
          },
          end: {
            line: 28,
            column: 9
          }
        }],
        line: 16
      },
      "2": {
        loc: {
          start: {
            line: 16,
            column: 12
          },
          end: {
            line: 16,
            column: 47
          }
        },
        type: "binary-expr",
        locations: [{
          start: {
            line: 16,
            column: 12
          },
          end: {
            line: 16,
            column: 26
          }
        }, {
          start: {
            line: 16,
            column: 30
          },
          end: {
            line: 16,
            column: 47
          }
        }],
        line: 16
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
exports.credentials = credentials;

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _fsExtra = _interopRequireDefault(require("fs-extra"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

async function credentials(config, _this) {
  cov_1iua2gc7za.f[0]++;
  cov_1iua2gc7za.s[0]++;

  try {
    cov_1iua2gc7za.s[1]++;

    if (!(await _fsExtra.default.pathExists(config))) {
      cov_1iua2gc7za.b[0][0]++;
      cov_1iua2gc7za.s[2]++;
      await _fsExtra.default.outputJson(config, {
        apiKey: '',
        apiSecret: ''
      });
    } else {
      cov_1iua2gc7za.b[0][1]++;
    }

    const {
      apiKey,
      apiSecret
    } = (cov_1iua2gc7za.s[3]++, await _fsExtra.default.readJson(config));
    cov_1iua2gc7za.s[4]++;

    if ((cov_1iua2gc7za.b[2][0]++, !apiKey.length) || (cov_1iua2gc7za.b[2][1]++, !apiSecret.length)) {
      cov_1iua2gc7za.b[1][0]++;
      cov_1iua2gc7za.s[5]++;

      _this.log(_chalk.default.red(`Credentials not found. Run ${_chalk.default.bold('chat config:set')} to generate a configuration file. ${_nodeEmoji.default.get('pensive')}`));

      cov_1iua2gc7za.s[6]++;

      _this.exit(0);
    } else {
      cov_1iua2gc7za.b[1][1]++;
    }

    cov_1iua2gc7za.s[7]++;
    return {
      apiKey,
      apiSecret
    };
  } catch (err) {
    cov_1iua2gc7za.s[8]++;

    _this.error(err);

    cov_1iua2gc7za.s[9]++;

    _this.exit(1);
  }
}