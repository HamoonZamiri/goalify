import { z } from "zod";

/**
 * Level entity schema
 */
export const LevelSchema = z.object({
	id: z.number(),
	level_up_xp: z.number(),
	cash_reward: z.number(),
});
export type Level = z.infer<typeof LevelSchema>;
