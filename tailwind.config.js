/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./frontend/pages/*.html"],
  theme: {
    extend: {
      boxShadow: {
        'customShadow': '0 0 1rem 0 rgba(0, 0, 0, 0.2)',
        'whiteShadow': '0 0 0 0.5px rgba(255, 255, 255, 0.4)'
      }
    },
    fontFamily:{
      'main': ['"Inter"'],
      'weird': ['"Handjet"',]
    }
  },
  plugins: [],
}

