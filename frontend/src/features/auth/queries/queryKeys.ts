/**
 * Query key factory for auth-related queries
 */
export const authKeys = {
	all: ["auth"] as const,
	sessions: () => [...authKeys.all, "sessions"] as const,
	currentUser: () => [...authKeys.sessions(), "current"] as const,
};
