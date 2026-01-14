/** @type {import('tailwindcss').Config} */
export default {
  darkMode: "class", // Habilitar dark mode basado en clase
  content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  theme: {
    extend: {
      screens: {
        // Breakpoints alineados con design tokens
        sm: "640px",
        md: "768px",
        lg: "1024px",
        xl: "1280px",
        "2xl": "1536px",
      },
    },
  },
  plugins: [],
};
