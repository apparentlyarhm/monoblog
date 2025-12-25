// @ts-check
import { defineConfig } from 'astro/config';

// https://astro.build/config
export default defineConfig({
    markdown: {
        shikiConfig: {
            theme: 'github-dark', // we will what all we can add later
            wrap: true,
        },
    }
});
