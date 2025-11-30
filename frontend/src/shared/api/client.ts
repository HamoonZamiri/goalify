import type { z } from "zod";
import useAuth from "@/shared/hooks/auth/useAuth";
import router from "@/router";
import type { ErrorResponse } from "@/shared/schemas/server-response.schema";
import { UserSchema } from "@/features/auth/schemas";
import { API_BASE, http } from "@/utils/constants";

let refreshPromise: Promise<unknown> | undefined;

/**
 * Refreshes the user's access token using the refresh token
 */
async function refreshUserToken() {
	if (refreshPromise) {
		return refreshPromise;
	}

	refreshPromise = (async () => {
		try {
			const { getUser, setUser, logout } = useAuth();
			const user = getUser();
			if (!user) return undefined;

			const res = await fetch(`${API_BASE}/users/refresh`, {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify({
					user_id: user.id,
					refresh_token: user.refresh_token,
				}),
			});

			const json: unknown = await res.json();
			if (!res.ok) {
				logout();
				router.push({ name: "Login" });
				return undefined;
			}

			const parsedJson = UserSchema.parse(json);
			setUser(parsedJson);
			return parsedJson;
		} finally {
			refreshPromise = undefined;
		}
	})();

	return refreshPromise;
}

/**
 * Type-safe fetch wrapper with Zod validation and automatic token refresh
 * Automatically adds Authorization header with access token from useAuth
 * @param url - The URL to fetch
 * @param schema - Zod schema to validate the response
 * @param options - Fetch options
 * @returns Parsed response or error
 */
export async function zodFetch<T>(
	url: string,
	schema: z.Schema<T>,
	options?: RequestInit,
): Promise<T | ErrorResponse> {
	const { getUser } = useAuth();
	const user = getUser();
	const accessToken = user?.access_token;

	// Automatically add auth headers if token exists
	const headers: HeadersInit = {
		"Content-Type": "application/json",
		...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
		...options?.headers,
	};

	const fetchOptions: RequestInit = {
		...options,
		headers,
	};

	let res = await fetch(url, fetchOptions);

	// Handle 204 No Content - no body to parse
	if (res.status === http.StatusNoContent) {
		return {} as T;
	}

	let json = await res.json();

	if (res.status === http.StatusUnauthorized) {
		const refreshedUser = await refreshUserToken();
		if (!refreshedUser) return json;

		res = await fetch(url, {
			...fetchOptions,
			headers: {
				...fetchOptions.headers,
				Authorization: `Bearer ${(refreshedUser as { access_token: string }).access_token}`,
			},
		});

		// Handle 204 No Content after retry
		if (res.status === http.StatusNoContent) {
			return {} as T;
		}

		json = await res.json();
	}

	if (!res.ok) {
		const error = json as ErrorResponse;
		error.statusCode = res.status;
		return error;
	}

	const parsedResponse = schema.parse(json);
	return parsedResponse;
}

/**
 * Gets the current user's access token
 * @deprecated Use zodFetch which automatically handles auth - this is kept for backward compatibility
 */
export function getAccessToken(): string | undefined {
	const { getUser } = useAuth();
	const user = getUser();
	return user?.access_token;
}

/**
 * Creates standard headers with authorization
 * @deprecated Use zodFetch which automatically handles auth - this is kept for backward compatibility
 */
export function createAuthHeaders(token?: string): HeadersInit {
	const headers: HeadersInit = {
		"Content-Type": "application/json",
	};

	const authToken = token || getAccessToken();
	if (authToken) {
		headers.Authorization = `Bearer ${authToken}`;
	}

	return headers;
}
