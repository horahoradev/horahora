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
  },
  typescript: {
    // !! WARN !!
    // Dangerously allow production builds to successfully complete even if
    // your project has type errors.
    // !! WARN !!
    ignoreBuildErrors: true,
  },
}

module.exports = nextJSConfig;
