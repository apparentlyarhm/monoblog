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

		category: z.string().optional(),
	}),
});

export const SITE = {
    title: "Arhm's blog",
    description: "I've now started writing about whatever I feel is interesting wiht no plans whatsoever",
    defaultImage: "/bg.png", 
};

export const collections = { blog };