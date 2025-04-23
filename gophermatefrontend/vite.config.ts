import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
  base: '/',
  plugins: [tailwindcss(), reactRouter(), tsconfigPaths(), react()],
  server: {
    port: 3000,
    historyApiFallback: true, // Add this line to handle SPA routing
  },
  esbuild: {
    jsx: 'automatic',
  },
});
