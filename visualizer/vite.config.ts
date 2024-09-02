import { defineConfig } from "vite";

export default defineConfig({
  build: {
    sourcemap: true,
    outDir: "dist",
    minify: "terser",
  },
});
