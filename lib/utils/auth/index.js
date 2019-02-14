"use strict";

var cov_e3hqmeehz = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/utils/auth/index.js";
  var hash = "4246f4a2b0906fb346bb877211353fcc919a39f8";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/utils/auth/index.js",
    statementMap: {
      "0": {
        start: {
          line: 6,
          column: 4
        },
        end: {
          line: 17,
          column: 5
        }
      },
      "1": {
        start: {
          line: 7,
          column: 38
        },
        end: {
          line: 7,
          column: 63
        }
      },
      "2": {
        start: {
          line: 8,
          column: 8
        },
        end: {
          line: 10,
          column: 9
        }
      },
      "3": {
        start: {
          line: 9,
          column: 12
        },
        end: {
          line: 9,
          column: 65
        }
      },
      "4": {
        start: {
          line: 12,
          column: 23
        },
        end: {
          line: 12,
          column: 56
        }
      },
      "5": {
        start: {
          line: 14,
          column: 8
        },
        end: {
          line: 14,
          column: 22
        }
      },
      "6": {
        start: {
          line: 16,
          column: 8
        },
        end: {
          line: 16,
          column: 38
        }
      }
    },
    fnMap: {
      "0": {
        name: "auth",
        decl: {
          start: {
            line: 5,
            column: 22
          },
          end: {
            line: 5,
            column: 26
          }
        },
        loc: {
          start: {
            line: 5,
            column: 42
          },
          end: {
            line: 18,
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
            line: 8,
            column: 8
          },
          end: {
            line: 10,
            column: 9
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 8,
            column: 8
          },
          end: {
            line: 10,
            column: 9
          }
        }, {
          start: {
            line: 8,
            column: 8
          },
          end: {
            line: 10,
            column: 9
          }
        }],
        line: 8
      },
      "1": {
        loc: {
          start: {
            line: 8,
            column: 12
          },
          end: {
            line: 8,
            column: 33
          }
        },
        type: "binary-expr",
        locations: [{
          start: {
            line: 8,
            column: 12
          },
          end: {
            line: 8,
            column: 19
          }
        }, {
          start: {
            line: 8,
            column: 23
          },
          end: {
            line: 8,
            column: 33
          }
        }],
        line: 8
      }
    },
    s: {
      "0": 0,
      "1": 0,
      "2": 0,
      "3": 0,
      "4": 0,
      "5": 0,
      "6": 0
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
exports.auth = auth;

var _streamChat = require("stream-chat");

var _config = require("../../utils/config");

async function auth(config, _this) {
  cov_e3hqmeehz.f[0]++;
  cov_e3hqmeehz.s[0]++;

  try {
    const {
      apiKey,
      apiSecret
    } = (cov_e3hqmeehz.s[1]++, await (0, _config.credentials)(config));
    cov_e3hqmeehz.s[2]++;

    if ((cov_e3hqmeehz.b[1][0]++, !apiKey) || (cov_e3hqmeehz.b[1][1]++, !apiSecret)) {
      cov_e3hqmeehz.b[0][0]++;
      cov_e3hqmeehz.s[3]++;
      return _this.error(`Missing config...`, {
        exit: 1
      });
    } else {
      cov_e3hqmeehz.b[0][1]++;
    }

    const client = (cov_e3hqmeehz.s[4]++, new _streamChat.StreamChat(apiKey, apiSecret));
    cov_e3hqmeehz.s[5]++;
    return client;
  } catch (err) {
    cov_e3hqmeehz.s[6]++;

    _this.error(err, {
      exit: 1
    });
  }
}