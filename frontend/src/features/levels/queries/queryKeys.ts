/**
 * Query key factory for level-related queries
 */
export const levelKeys = {
	all: ["levels"] as const,
	details: () => [...levelKeys.all, "detail"] as const,
	detail: (id: number) => [...levelKeys.details(), id] as const,
};
