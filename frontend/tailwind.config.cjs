/** @type {import('tailwindcss').Config} */

const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      fontFamily: {
        'sans': ['Manrope', defaultTheme.fontFamily.sans],
      },
      animation: {
        'spin-fast': 'spin .75s linear infinite',
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
