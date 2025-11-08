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

export const LevelByIdParams = z.object({
	levelId: z.number(),
});
export type LevelByIdParams = z.infer<typeof LevelByIdParams>;
