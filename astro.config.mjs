// @ts-check
import { defineConfig } from 'astro/config';

import sitemap from '@astrojs/sitemap';

// https://astro.build/config
export default defineConfig({
  markdown: {
      shikiConfig: {
          theme: 'github-dark', // we will what all we can add later
          wrap: true,
      },
  },

  site: 'https://blog.arhm.dev',
  integrations: [sitemap()],
});