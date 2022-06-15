// craco.config.js
/**
 * @type {import("@craco/craco").CracoConfig}
 */
module.exports = {
    style: {
        postcssOptions: {
            plugins: [
                require('tailwindcss'),
                require('autoprefixer'),
            ],
        },
    }
}