import colors from 'tailwindcss/colors'

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
        // The customer site's brand red. 500 is the bright accent; buttons use
        // 600 so white text stays readable. 700 is the logo's red.
        primary: {
          DEFAULT: '#c3241a',
          50: '#fff2f1',
          100: '#ffe1de',
          200: '#fec7c1',
          300: '#fb9e94',
          400: '#ef5350',
          500: '#f44336',
          600: '#c3241a',
          700: '#9b1414',
          800: '#7f1513',
          900: '#681614',
        },
        // The logo's leaf green, for "live", "paid" and success.
        leaf: {
          50: '#f4fbea',
          100: '#e6f5d0',
          400: '#8cc343',
          500: '#72a82f',
          600: '#578624',
          700: '#43671f',
        },
        // Greys are the slate blues the customer site uses, so dark mode is
        // the same navy across both apps.
        gray: colors.slate,
        secondary: {
          DEFAULT: '#FFFFFF',
          dark: '#F8F9FA',
        },
        accent: {
          DEFAULT: '#64748b',
          light: '#f8fafc',
          dark: '#1e293b',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        display: ['Poppins', 'Inter', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        soft: '0 1px 2px rgb(15 23 42 / 0.04), 0 4px 16px rgb(15 23 42 / 0.06)',
        lift: '0 2px 4px rgb(15 23 42 / 0.06), 0 12px 32px rgb(15 23 42 / 0.14)',
      },
      keyframes: {
        'slide-up': { '0%': { transform: 'translateY(100%)' }, '100%': { transform: 'translateY(0)' } },
        'slide-left': { '0%': { transform: 'translateX(100%)' }, '100%': { transform: 'translateX(0)' } },
        'fade-in': { '0%': { opacity: '0' }, '100%': { opacity: '1' } },
        pop: { '0%': { transform: 'scale(1)' }, '40%': { transform: 'scale(1.3)' }, '100%': { transform: 'scale(1)' } },
      },
      animation: {
        'slide-up': 'slide-up 0.28s cubic-bezier(0.32, 0.72, 0, 1)',
        'slide-left': 'slide-left 0.28s cubic-bezier(0.32, 0.72, 0, 1)',
        'fade-in': 'fade-in 0.2s ease-out',
        pop: 'pop 0.35s ease-out',
      },
    },
  },
  plugins: [],
}
