import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import fs from 'fs';

export default defineConfig({
  resolve: {
    alias: {
      $src: "/src",
    },
  },
	plugins: [sveltekit()],
    server: {
      host: true,
      port: 8080
    }
});
