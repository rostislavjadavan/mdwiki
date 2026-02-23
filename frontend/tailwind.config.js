/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        'gh-bg': '#ffffff',
        'gh-border': '#d0d7de',
        'gh-header': '#24292f',
        'gh-link': '#0969da',
        'gh-btn-primary': '#2da44e',
        'gh-btn-danger': '#cf222e',
        'gh-text': '#1f2328',
        'gh-muted': '#656d76',
      },
    },
  },
  plugins: [],
}
