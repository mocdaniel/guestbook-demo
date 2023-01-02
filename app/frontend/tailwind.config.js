/** @type {import('tailwindcss').Config} */
const colors = require("tailwindcss/colors");

module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx, vue, svg}",
  ],
  theme: {
    extend: {
      colors: {
        brand: "#4D9BF8",
        dark: "#070F1F",
        mid: "#B8BFC6",
        light: "#FFFFFF",
        grid: "#111828"
      },
    },
  },
}
