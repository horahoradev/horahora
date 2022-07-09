// @ts-check
const path = require("path");

/**
 * @type {import("next").NextConfig}
 */
const nextJSConfig = {
  reactStrictMode: true,
  swcMinify: true,
  eslint: {
    dirs: ["environment", "src"],
  },
  sassOptions: {
    includePaths: [path.join(__dirname, "src", "styles")],
  }
};

module.exports = nextJSConfig;
