/**
 * @type {Config}
 */
const postCSSConfig = {
  plugins:
    process.env.NODE_ENV === "production"
      ? {
          "tailwindcss": {},
          "postcss-flexbugs-fixes": {},
          "postcss-preset-env": {
            autoprefixer: {
              flexbox: "no-2009",
            },
            stage: 3,
            features: {
              "custom-properties": false,
            },
          },
        }
      : {
          "tailwindcss": {},
        },
};

module.exports = postCSSConfig;
