import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import path from "path";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  define: {
    "process.env": process.env,
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
      "@go-kickstart/zod": path.resolve(__dirname, "../../packages/zod/src"),
    },
  },
});
