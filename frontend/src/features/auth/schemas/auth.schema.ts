import { z } from "zod";
import { createServerResponseSchema } from "@/shared/schemas/server-response.schema";

/**
 * User entity schema
 */
export const UserSchema = z.object({
	id: z.string().uuid(),
	email: z.string().email(),
	xp: z.number(),
	level_id: z.number(),
	cash_available: z.number(),
	access_token: z.string(),
	refresh_token: z.string().uuid(),
});
export type User = z.infer<typeof UserSchema>;

/**
 * API response schemas
 */
export const UserResponseSchema = createServerResponseSchema(UserSchema);
