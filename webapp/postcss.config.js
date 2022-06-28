const IS_PRODUCTION = process.env.NODE_ENV === "production";

const postCSSConfig = {
  plugins: {
    tailwindcss: {},
  },
};

if (IS_PRODUCTION) {
  const prodPlugins = {
    "postcss-flexbugs-fixes": {},
    "postcss-preset-env": {
      autoprefixer: {
        flexbox: "no-2009",
        grid: "autoplace"
      },
      stage: 3,
      features: {
        "custom-properties": false,
      },
    },
  };

  postCSSConfig.plugins = { ...postCSSConfig.plugins, ...prodPlugins };
}

module.exports = postCSSConfig;
