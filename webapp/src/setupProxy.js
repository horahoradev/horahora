// create-react-app proxy for API (to avoid CORS errors)

const { createProxyMiddleware } = require("http-proxy-middleware");

module.exports = function (app) {
  // app.use(
  //   "/api",
  //   createProxyMiddleware({
  //     target: "http://frontapi:8083",
  //     changeOrigin: true,
  //     pathRewrite: { "^/api": "" },
  //   })
  // );
  // app.use(
  //   "/static/images",
  //   createProxyMiddleware({
  //     target: "http://nginx:80",
  //     changeOrigin: true,
  //     pathRewrite: { "^/static/images": "" },
  //   })
  // );
};
