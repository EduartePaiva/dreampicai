/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./view/**/*.templ", "./**/*.templ"],
  safelist: [],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["dark"],
  },
};
