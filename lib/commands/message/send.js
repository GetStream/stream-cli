"use strict";

var cov_1cp42l8jk6 = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/commands/message/send.js";
  var hash = "63c9f962c2c4cc4f99ee6f590be8fc05f0d330c4";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/commands/message/send.js",
    statementMap: {
      "0": {
        start: {
          line: 42,
          column: 26
        },
        end: {
          line: 42,
          column: 49
        }
      },
      "1": {
        start: {
          line: 44,
          column: 8
        },
        end: {
          line: 80,
          column: 9
        }
      },
      "2": {
        start: {
          line: 45,
          column: 27
        },
        end: {
          line: 48,
          column: 13
        }
      },
      "3": {
        start: {
          line: 50,
          column: 12
        },
        end: {
          line: 53,
          column: 15
        }
      },
      "4": {
        start: {
          line: 55,
          column: 12
        },
        end: {
          line: 55,
          column: 73
        }
      },
      "5": {
        start: {
          line: 56,
          column: 28
        },
        end: {
          line: 56,
          column: 69
        }
      },
      "6": {
        start: {
          line: 58,
          column: 28
        },
        end: {
          line: 60,
          column: 13
        }
      },
      "7": {
        start: {
          line: 62,
          column: 12
        },
        end: {
          line: 64,
          column: 13
        }
      },
      "8": {
        start: {
          line: 63,
          column: 16
        },
        end: {
          line: 63,
          column: 68
        }
      },
      "9": {
        start: {
          line: 66,
          column: 12
        },
        end: {
          line: 66,
          column: 47
        }
      },
      "10": {
        start: {
          line: 68,
          column: 28
        },
        end: {
          line: 74,
          column: 13
        }
      },
      "11": {
        start: {
          line: 76,
          column: 12
        },
        end: {
          line: 76,
          column: 50
        }
      },
      "12": {
        start: {
          line: 77,
          column: 12
        },
        end: {
          line: 77,
          column: 24
        }
      },
      "13": {
        start: {
          line: 79,
          column: 12
        },
        end: {
          line: 79,
          column: 41
        }
      },
      "14": {
        start: {
          line: 84,
          column: 0
        },
        end: {
          line: 84,
          column: 56
        }
      }
    },
    fnMap: {
      "0": {
        name: "(anonymous_0)",
        decl: {
          start: {
            line: 41,
            column: 4
          },
          end: {
            line: 41,
            column: 5
          }
        },
        loc: {
          start: {
            line: 41,
            column: 16
          },
          end: {
            line: 81,
            column: 5
          }
        },
        line: 41
      }
    },
    branchMap: {
      "0": {
        loc: {
          start: {
            line: 62,
            column: 12
          },
          end: {
            line: 64,
            column: 13
          }
        },
        type: "if",
        locations: [{
          start: {
            line: 62,
            column: 12
          },
          end: {
            line: 64,
            column: 13
          }
        }, {
          start: {
            line: 62,
            column: 12
          },
          end: {
            line: 64,
            column: 13
          }
        }],
        line: 62
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
      "14": 0
    },
    f: {
      "0": 0
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
exports.MessageSend = void 0;

var _command = require("@oclif/command");

var _nodeEmoji = _interopRequireDefault(require("node-emoji"));

var _chalk = _interopRequireDefault(require("chalk"));

var _path = _interopRequireDefault(require("path"));

var _uuid = _interopRequireDefault(require("uuid"));

var _auth = require("../../utils/auth");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class MessageSend extends _command.Command {
  async run() {
    cov_1cp42l8jk6.f[0]++;
    const {
      flags
    } = (cov_1cp42l8jk6.s[0]++, this.parse(MessageSend));
    cov_1cp42l8jk6.s[1]++;

    try {
      const client = (cov_1cp42l8jk6.s[2]++, await (0, _auth.auth)(_path.default.join(this.config.configDir, 'config.json'), this));
      cov_1cp42l8jk6.s[3]++;
      await client.updateUser({
        id: flags.uid,
        role: 'admin'
      });
      cov_1cp42l8jk6.s[4]++;
      await client.setUser({
        id: flags.uid,
        status: 'invisible'
      });
      const channel = (cov_1cp42l8jk6.s[5]++, client.channel(flags.type, flags.channel));
      const payload = (cov_1cp42l8jk6.s[6]++, {
        text: flags.message
      });
      cov_1cp42l8jk6.s[7]++;

      if (flags.attachments) {
        cov_1cp42l8jk6.b[0][0]++;
        cov_1cp42l8jk6.s[8]++;
        payload.attachments = JSON.parse(flags.attachments);
      } else {
        cov_1cp42l8jk6.b[0][1]++;
      }

      cov_1cp42l8jk6.s[9]++;
      await channel.sendMessage(payload);
      const message = (cov_1cp42l8jk6.s[10]++, _chalk.default.blue(`Message ${_chalk.default.bold(flags.message)} has been sent to the ${_chalk.default.bold(flags.channel)} channel by ${_chalk.default.bold(flags.uid)}!`));
      cov_1cp42l8jk6.s[11]++;
      this.log(message, _nodeEmoji.default.get('smile'));
      cov_1cp42l8jk6.s[12]++;
      this.exit();
    } catch (err) {
      cov_1cp42l8jk6.s[13]++;
      this.error(err, {
        exit: 1
      });
    }
  }

}

exports.MessageSend = MessageSend;

_defineProperty(MessageSend, "flags", {
  id: _command.flags.string({
    char: 'i',
    description: _chalk.default.blue.bold('Channel ID.'),
    default: (0, _uuid.default)(),
    required: false
  }),
  user: _command.flags.string({
    char: 'u',
    description: _chalk.default.blue.bold('ID of user.'),
    default: '*',
    required: true
  }),
  type: _command.flags.string({
    char: 't',
    description: _chalk.default.blue.bold('Type of channel.'),
    options: ['livestream', 'messaging', 'gaming', 'commerce', 'team'],
    required: true
  }),
  message: _command.flags.string({
    char: 'm',
    description: _chalk.default.blue.bold('Message to send.'),
    required: true
  }),
  attachments: _command.flags.string({
    char: 'a',
    description: _chalk.default.blue.bold('JSON payload of attachments'),
    required: false
  })
});

cov_1cp42l8jk6.s[14]++;
MessageSend.description = 'Send messages to a channel.';