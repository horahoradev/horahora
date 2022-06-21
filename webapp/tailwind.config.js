/** @type {import('tailwindcss').Config} */
const tailwindConfig = {
  content: ["./public/**/*.html", "./src/**/*.{js,jsx,ts,tsx}"],
  important: true,
  theme: {
    extend: {},
  },
  plugins: [],
};

module.exports = tailwindConfig;
