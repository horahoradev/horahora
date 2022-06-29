// @ts-check

/**
 * @type {import("next").NextConfig}
 */
const nextJSConfig = {
  reactStrictMode: true,
  swcMinify: true,
  eslint: {
    dirs: ["environment", "src"]
  },
  experimental: {
    newNextLinkBehavior: true
  }
};

module.exports = nextJSConfig;
