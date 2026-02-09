/** @type {import('tailwindcss').Config} */
export default {
    darkMode: 'class',
    content: [
      "./index.html",
      "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
      extend: {
        colors: {
          primary: {
            DEFAULT: '#E40A2D',
            50: '#FEF2F3',
            100: '#FDE5E8',
            200: '#FABFC7',
            300: '#F799A6',
            400: '#F24C66',
            500: '#E40A2D',
            600: '#CD0929',
            700: '#89061B',
            800: '#670514',
            900: '#44030E',
          },
          secondary: {
            DEFAULT: '#FFFFFF',
            dark: '#F8F9FA',
          },
          accent: {
            DEFAULT: '#6C757D',
            light: '#F8F9FA',
            dark: '#343A40',
          }
        },
        fontFamily: {
          sans: ['Inter', 'system-ui', 'sans-serif'],
        },
      },
    },
    plugins: [],
  }