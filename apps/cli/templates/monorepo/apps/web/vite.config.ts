import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import path from "path";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  test: {
    css: true,
    environment: "jsdom",
    setupFiles: "./src/testing/setup.ts",
    globals: true,
    clearMocks: true,
    restoreMocks: true,
    env: {
      VITE_API_URL: "http://localhost:3001",
      VITE_ENV: "development",
      VITE_GOOGLE_AUTH_ENABLED: "false",
    },
  },
  server: {
    port: 3000,
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@go-kickstart/openapi": path.resolve(
        __dirname,
        "../../packages/openapi/src",
      ),
      "@go-kickstart/ui": path.resolve(__dirname, "../../packages/ui/src"),
      "@go-kickstart/zod": path.resolve(__dirname, "../../packages/zod/src"),
    },
  },
  optimizeDeps: {
    exclude: ["@go-kickstart/openapi", "@go-kickstart/ui", "@go-kickstart/zod"],
  },
});
