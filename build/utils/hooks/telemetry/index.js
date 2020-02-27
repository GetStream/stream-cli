"use strict";

var _md = _interopRequireDefault(require("md5"));

var _rollbar = _interopRequireDefault(require("rollbar"));

var _config = require("../../config");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

module.exports = async ({
  ctx,
  error
}) => {
  const {
    name,
    email,
    apiKey,
    apiSecret,
    apiBaseUrl,
    environment,
    telemetry
  } = await (0, _config.credentials)(ctx);

  if (telemetry === 'true') {
    const rollbar = new _rollbar.default({
      accessToken: '4ba9cceb33fe4543b7a114fc354b870c',
      captureUncaught: true,
      captureUnhandledRejections: true,
      environment: 'production'
    });
    rollbar.error(error, {
      person: {
        id: (0, _md.default)(email),
        username: email,
        name,
        email
      },
      api: {
        key: apiKey,
        secret: apiSecret,
        url: apiBaseUrl
      },
      environment
    });
  }

  if (typeof error === 'object' && (error.code !== 'EEXIT' || error.code === undefined)) {
    ctx.error(error.message, {
      exit: 1
    });
  }
};