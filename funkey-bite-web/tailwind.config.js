/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: [
    './components/**/*.{js,vue,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './plugins/**/*.{js,ts}',
    './composables/**/*.{js,ts}',
    './utils/**/*.{js,ts}',
    './nuxt.config.{js,ts}',
    './app.vue'
  ],
  theme: {
    extend: {
      colors: {
        // Brand red, from the logo (#9b1414 is 700): buttons, prices, active states.
        brand: {
          50: '#fff2f1',
          100: '#ffe1de',
          200: '#fec7c1',
          300: '#fb9e94',
          400: '#f0695b',
          500: '#e03b2d',
          600: '#c3241a',
          700: '#9b1414',
          800: '#7f1513',
          900: '#681614',
        },
        // Leaf green, from the logo's spoon (#8cc343 is 400): open, fresh, success.
        leaf: {
          50: '#f4fbea',
          100: '#e6f5d0',
          200: '#cdeaa3',
          300: '#addb6f',
          400: '#8cc343',
          500: '#72a82f',
          600: '#578624',
          700: '#43671f',
          800: '#37521d',
        },
        // Warm amber: highlights and the gradient's second stop.
        accent: {
          50: '#fffbeb',
          100: '#fef3c7',
          200: '#fde68a',
          300: '#fcd34d',
          400: '#fbbf24',
          500: '#f59e0b',
          600: '#d97706',
          700: '#b45309',
        },
        // Warm neutrals for text and surfaces, so the page doesn't look grey.
        ink: {
          50: '#faf8f6',
          100: '#f3efeb',
          200: '#e7e0d9',
          300: '#d3c9bf',
          400: '#a59a8f',
          500: '#7a6f66',
          600: '#5c524a',
          700: '#443c36',
          800: '#2a2420',
          900: '#1a1512',
          950: '#110d0b',
        },
      },
      fontFamily: {
        display: ['Poppins', 'system-ui', 'sans-serif'],
        body: ['Inter', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        soft: '0 1px 2px rgb(26 21 18 / 0.04), 0 4px 16px rgb(26 21 18 / 0.06)',
        lift: '0 2px 4px rgb(26 21 18 / 0.06), 0 12px 32px rgb(26 21 18 / 0.12)',
        glow: '0 8px 24px -8px rgb(195 36 26 / 0.55)',
      },
      spacing: {
        // Space for the bottom tab bar on phones, including the home indicator.
        tabbar: 'calc(4rem + env(safe-area-inset-bottom))',
      },
      keyframes: {
        'sheet-up': {
          '0%': { transform: 'translateY(100%)' },
          '100%': { transform: 'translateY(0)' },
        },
        'fade-in': {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        'rise-in': {
          '0%': { opacity: '0', transform: 'translateY(12px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        pop: {
          '0%': { transform: 'scale(1)' },
          '40%': { transform: 'scale(1.25)' },
          '100%': { transform: 'scale(1)' },
        },
        drift: {
          '0%, 100%': { transform: 'translate3d(0, 0, 0) rotate(0deg)' },
          '50%': { transform: 'translate3d(0, -8px, 0) rotate(2deg)' },
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
      },
      animation: {
        'sheet-up': 'sheet-up 0.28s cubic-bezier(0.32, 0.72, 0, 1)',
        'fade-in': 'fade-in 0.2s ease-out',
        'rise-in': 'rise-in 0.5s cubic-bezier(0.22, 1, 0.36, 1) both',
        pop: 'pop 0.35s ease-out',
        drift: 'drift 6s ease-in-out infinite',
        shimmer: 'shimmer 1.6s linear infinite',
      },
    },
  },
  plugins: [],
}
