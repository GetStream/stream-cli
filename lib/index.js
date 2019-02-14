"use strict";

var cov_260iv6tdkt = function () {
  var path = "/Users/parsons/Code/stream/stream-cli/src/index.js";
  var hash = "4f0b1c05824df3449daa1768f80b3204477d9a6b";

  var Function = function () {}.constructor;

  var global = new Function("return this")();
  var gcv = "__coverage__";
  var coverageData = {
    path: "/Users/parsons/Code/stream/stream-cli/src/index.js",
    statementMap: {},
    fnMap: {},
    branchMap: {},
    s: {},
    f: {},
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
exports.default = void 0;

var _command = require("@oclif/command");

var _default = _command.Command;
exports.default = _default;