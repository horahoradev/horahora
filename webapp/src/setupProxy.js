// create-react-app proxy for API (to avoid CORS errors)

const { createProxyMiddleware } = require("http-proxy-middleware");

module.exports = function (app) {
  app.use(
    "/api",
    createProxyMiddleware({
      target: "http://localhost:8083",
      changeOrigin: true,
      pathRewrite: { "^/api": "" },
    })
  );
};
