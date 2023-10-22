/** @type {import('tailwindcss').Config} */
const colors = require('tailwindcss/colors');
module.exports = {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    // colors: {
    //   black: colors.black,
    //   'dark-purple': '#52057B',
    //   purple: '#892CDC',
    //   'light-purple': '#BC6FF1'
    // },
    extend: {
      fontFamily: {
        ubuntu: ['UbuntuMono'],
        icomoon: ['icomoon']
      },
      colors: {
        black: colors.black,
        'dark-purple': '#52057B',
        purple: '#892CDC',
        'light-purple': '#BC6FF1'
      },
    },
  },
  plugins: [],
}
