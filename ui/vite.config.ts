import { defineConfig } from 'vite';
import { resolve } from 'path';

export default defineConfig({
  build: {
    lib: {
      entry: resolve(__dirname, 'src/index.ts'),
      name: 'DssUI',
      fileName: 'dss-ui',
      formats: ['es', 'iife'],
    },
    rollupOptions: {
      output: {
        // Ensure CSS is inlined in the JS bundle
        assetFileNames: 'dss-ui.[ext]',
      },
    },
    minify: 'esbuild',
    sourcemap: true,
  },
  define: {
    'process.env.NODE_ENV': JSON.stringify('production'),
  },
});
