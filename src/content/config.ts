import { defineCollection, z } from 'astro:content';

const blog = defineCollection({
	schema: z.object({
		title: z.string(),
		date: z.string(),
		tags: z.string().optional(),
		description: z.string().optional(),
		image: z.string().optional(),
	}),
});

export const SITE = {
    title: "arhm's blog",
    description: "a site where arhm writes BS",
    defaultImage: "/bg.png", 
};

export const collections = { blog };