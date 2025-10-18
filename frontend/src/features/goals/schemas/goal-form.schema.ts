import { z } from "zod";

/**
 * Create goal form validation schema
 */
export const createGoalFormSchema = z.object({
	title: z
		.string()
		.min(1, "Title is required")
		.max(100, "Title must be 100 characters or less"),
	description: z
		.string()
		.min(1, "Description is required")
		.max(500, "Description must be 500 characters or less"),
	category_id: z.string().uuid("Invalid category ID"),
});

/**
 * Create goal category form validation schema
 */
export const createGoalCategoryFormSchema = z.object({
	title: z
		.string()
		.min(1, "Title is required")
		.max(50, "Title must be 50 characters or less"),
	xp_per_goal: z
		.number()
		.min(1, "XP must be at least 1")
		.max(1000, "XP must be 1000 or less"),
});

/**
 * Update goal form validation schema
 */
export const updateGoalFormSchema = z.object({
	title: z.string().min(1).max(100).optional(),
	description: z.string().min(1).max(500).optional(),
	status: z.enum(["complete", "not_complete"]).optional(),
});

/**
 * Update goal category form validation schema
 */
export const updateGoalCategoryFormSchema = z.object({
	title: z.string().min(1).max(50).optional(),
	xp_per_goal: z.number().min(1).max(1000).optional(),
});

/**
 * Inferred form types
 */
export type CreateGoalFormData = z.infer<typeof createGoalFormSchema>;
export type CreateGoalCategoryFormData = z.infer<
	typeof createGoalCategoryFormSchema
>;
export type UpdateGoalFormData = z.infer<typeof updateGoalFormSchema>;
export type UpdateGoalCategoryFormData = z.infer<
	typeof updateGoalCategoryFormSchema
>;
