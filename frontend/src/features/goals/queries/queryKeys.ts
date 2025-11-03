/**
 * Query key factory for goals domain
 * Follows TanStack Query best practices for hierarchical key structure
 */

export const goalKeys = {
	all: ["goals"] as const,
	lists: () => [...goalKeys.all, "list"] as const,
	list: (filters?: Record<string, unknown>) =>
		[...goalKeys.lists(), filters] as const,
	details: () => [...goalKeys.all, "detail"] as const,
	detail: (id: string) => [...goalKeys.details(), id] as const,
};

export const categoryKeys = {
	all: ["categories"] as const,
	lists: () => [...categoryKeys.all, "list"] as const,
	list: (filters?: Record<string, unknown>) =>
		[...categoryKeys.lists(), filters] as const,
	details: () => [...categoryKeys.all, "detail"] as const,
	detail: (id: string) => [...categoryKeys.details(), id] as const,
};
