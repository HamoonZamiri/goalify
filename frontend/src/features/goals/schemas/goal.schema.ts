import { z } from "zod";
import {
	createServerResponseSchema,
	createArraySchema,
} from "@/shared/schemas/server-response.schema";

/**
 * Goal entity schema
 */
export const GoalSchema = z.object({
	id: z.string().uuid(),
	title: z.string(),
	description: z.string(),
	category_id: z.string().uuid(),
	user_id: z.string().uuid(),
	status: z.enum(["complete", "not_complete"]),
	created_at: z.string(),
	updated_at: z.string(),
});

/**
 * Goal Category entity schema
 */
export const GoalCategorySchema = z.object({
	id: z.string().uuid(),
	title: z.string(),
	xp_per_goal: z.number(),
	user_id: z.string().uuid(),
	goals: z.array(GoalSchema),
	created_at: z.string(),
	updated_at: z.string(),
});

/**
 * API response schemas
 */
export const GoalResponseSchema = createServerResponseSchema(GoalSchema);
export const GoalCategoryResponseSchema =
	createServerResponseSchema(GoalCategorySchema);
export const GoalArraySchema = createArraySchema(GoalSchema);
export const GoalCategoryArraySchema = createArraySchema(GoalCategorySchema);
export const GoalResponseArraySchema =
	createServerResponseSchema(GoalArraySchema);
export const GoalCategoryResponseArraySchema = createServerResponseSchema(
	GoalCategoryArraySchema,
);

/**
 * Inferred types from schemas
 */
export type Goal = z.infer<typeof GoalSchema>;
export type GoalCategory = z.infer<typeof GoalCategorySchema>;
