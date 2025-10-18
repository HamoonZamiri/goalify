import type { z } from "zod";
import { API_BASE, http } from "@/utils/constants";
import type { ErrorResponse } from "@/shared/schemas/server-response.schema";
import useAuth from "@/hooks/auth/useAuth";
import router from "@/router";
import { Schemas } from "@/utils/schemas";

let refreshPromise: Promise<unknown> | undefined = undefined;

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

			const parsedJson = Schemas.UserSchema.parse(json);
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
	let res = await fetch(url, options);
	let json = await res.json();

	if (res.status === http.StatusUnauthorized) {
		const refreshedUser = await refreshUserToken();
		if (!refreshedUser) return json;

		res = await fetch(url, {
			...options,
			headers: {
				...options?.headers,
				Authorization: `Bearer ${(refreshedUser as { access_token: string }).access_token}`,
			},
		});
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
 */
export function getAccessToken(): string | undefined {
	const { getUser } = useAuth();
	const user = getUser();
	return user?.access_token;
}

/**
 * Creates standard headers with authorization
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
