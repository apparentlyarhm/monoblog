import { defineCollection, z } from 'astro:content';

const blog = defineCollection({
	schema: z.object({
		// Blog's title
		title: z.string(),

		// Date of the blog post
		date: z.string(),

		// Optional tags for the blog post
		tags: z.string().optional(),

		// Optional description for the blog post
		description: z.string().optional(),

		// Optional image for the blog post
		image: z.string().optional(),
	}),
});

export const SITE = {
    title: "arhm's blog",
    description: "a site where arhm writes BS",
    defaultImage: "/bg.png", 
};

export const collections = { blog };