// craco.config.js
module.exports = {
    style: {
        postcssOptions: {
            plugins: [
                require('tailwindcss'),
                require('autoprefixer'),
            ],
        },
    },
}