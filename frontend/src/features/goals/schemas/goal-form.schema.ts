import { z } from "zod";

/**
 * Create goal form validation schema
 */
export const createGoalFormSchema = z.object({
	title: z
		.string()
		.min(1, "Title is required")
		.max(255, "Title must be 255 characters or less"),
	description: z
		.string()
		.min(1, "Description is required")
		.max(255, "Description must be 255 characters or less"),
	category_id: z.string().uuid("Invalid category ID"),
});
export type CreateGoalFormData = z.infer<typeof createGoalFormSchema>;

/**
 * Create goal category form validation schema
 */
export const createGoalCategoryFormSchema = z.object({
	title: z
		.string()
		.min(1, "Title is required")
		.max(255, "Title must be 255 characters or less"),
	xp_per_goal: z.coerce
		.number({ invalid_type_error: "XP must be a number" })
		.positive("XP must be positive")
		.min(1, "XP must be at least 1")
		.max(100, "XP must be 100 or less"),
});
export type CreateGoalCategoryFormData = z.infer<
	typeof createGoalCategoryFormSchema
>;

/**
 * Update goal form validation schema
 */
export const updateGoalFormSchema = z.object({
	title: z.string().min(1).max(255).optional(),
	description: z.string().min(1).max(255).optional(),
	status: z.enum(["complete", "not_complete"]).optional(),
});
export type UpdateGoalFormData = z.infer<typeof updateGoalFormSchema>;

/**
 * Update goal category form validation schema
 */
export const updateGoalCategoryFormSchema = z.object({
	title: z.string().min(1).max(255).optional(),
	xp_per_goal: z.number().min(1).max(100).optional(),
});
export type UpdateGoalCategoryFormData = z.infer<
	typeof updateGoalCategoryFormSchema
>;

/**
 * Edit goal form validation schema (for inline editing with all required fields)
 */
export const editGoalFormSchema = z.object({
	title: z
		.string()
		.min(1, "Title is required")
		.max(255, "Title must be 255 characters or less"),
	description: z
		.string()
		.min(1, "Description is required")
		.max(255, "Description must be 255 characters or less"),
	status: z.enum(["complete", "not_complete"]),
});
export type EditGoalFormData = z.infer<typeof editGoalFormSchema>;

/**
 * Edit goal category form validation schema (for inline editing with all required fields)
 */
export const editGoalCategoryFormSchema = z.object({
	title: z
		.string()
		.min(1, "Title is required")
		.max(255, "Title must be 255 characters or less"),
	xp_per_goal: z
		.number({ invalid_type_error: "XP must be a number" })
		.positive("XP must be positive")
		.min(1, "XP must be at least 1")
		.max(100, "XP must be 100 or less"),
});
export type EditGoalCategoryFormData = z.infer<
	typeof editGoalCategoryFormSchema
>;

export const ResetGoalCategoryParams = z.object({
	category_id: z.string().uuid(),
});
export type ResetGoalCategoryParams = z.infer<typeof ResetGoalCategoryParams>;
