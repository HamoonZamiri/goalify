import type { LevelByIdParams } from "../schemas";

/**
 * Query key factory for level-related queries
 */
export const levelKeys = {
	all: ["levels"] as const,
	details: () => [...levelKeys.all, "detail"] as const,
	detail: (params: LevelByIdParams) =>
		[...levelKeys.details(), params] as const,
};

export const getLevelByIdParams = (
	queryKey: ReturnType<typeof levelKeys.detail>,
): LevelByIdParams => {
	return queryKey[2];
};
