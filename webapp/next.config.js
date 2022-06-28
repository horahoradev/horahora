/**
 * @type {import("next").NextConfig}
 */
const nextJSConfig = {
  reactStrictMode: true,
  swcMinify: true,
  eslint: {
    dirs: ["environment", "src"]
  }
};

module.exports = nextJSConfig;
