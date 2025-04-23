import { defineConfig } from "vite";

export default defineConfig({
  build: {
    sourcemap: true,
    outDir: "dist",
    minify: "terser",
    rollupOptions: {
      output: {
        entryFileNames: `assets/[name].js`,
        chunkFileNames: `assets/[name].js`,
        assetFileNames: `assets/[name].[ext]`
      }
    }
  },
});
