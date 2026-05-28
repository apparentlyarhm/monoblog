// @ts-check
import { defineConfig } from 'astro/config';

import sitemap from '@astrojs/sitemap';
import rehypeKatex from 'rehype-katex';
import remarkMath from 'remark-math';

// https://astro.build/config
export default defineConfig({
  markdown: {
    remarkPlugins: [remarkMath],
    rehypePlugins: [rehypeKatex],
      shikiConfig: {
          theme: 'houston', // we will what all we can add later
          wrap: true,
      },
  },

  site: 'https://blog.arhm.dev',
  integrations: [sitemap()],
});