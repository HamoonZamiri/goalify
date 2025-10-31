import { QueryClient } from "@tanstack/vue-query";

/**
 * TanStack Query client configuration
 * Centralized configuration for caching, refetching, and error handling
 */
export const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			staleTime: 1000 * 60 * 5, // 5 minutes
			gcTime: 1000 * 60 * 10, // 10 minutes (formerly cacheTime)
			retry: 1,
			refetchOnWindowFocus: false,
			refetchOnReconnect: true,
		},
		mutations: {
			retry: 0,
		},
	},
});
