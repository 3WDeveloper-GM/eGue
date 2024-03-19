/** @type {import('tailwindcss').Config} */

const defaultTheme = require('tailwindcss/defaultTheme')


export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: "#2e3440",
        secondary: "#4c566a",
        third: "#d8dee9",
        fourth: "#eceff4",
        fifth: "#5e81ac"
      },
      fontFamily: {
        sans: ['Inter var', ...defaultTheme.fontFamily.sans],
      }
    },
  },
  plugins: [require("daisyui")],
}