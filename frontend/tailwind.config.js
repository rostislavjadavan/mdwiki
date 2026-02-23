/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        'gh-bg': 'var(--gh-bg)',
        'gh-border': 'var(--gh-border)',
        'gh-header': 'var(--gh-header)',
        'gh-link': 'var(--gh-link)',
        'gh-btn-primary': 'var(--gh-btn-primary)',
        'gh-btn-danger': 'var(--gh-btn-danger)',
        'gh-text': 'var(--gh-text)',
        'gh-muted': 'var(--gh-muted)',
        'gh-subtle': 'var(--gh-subtle)',
      },
    },
  },
  plugins: [],
}
