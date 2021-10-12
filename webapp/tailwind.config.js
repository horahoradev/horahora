module.exports = {
  important: true,
  mode: "jit",
  purge: [
    //
    "./public/**/*.html",
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {},
  },
  variants: {
    extend: {},
  },
  plugins: [],
};
