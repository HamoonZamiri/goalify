import authState from "@/state/auth";
import { API_BASE, http } from "./constants";
import { Schemas } from "./schemas";
import type { z } from "zod";

type ServerResponse<T> = {
  message: string;
  data: T;
};

async function refreshUserToken(): Promise<void | Error> {
  if (!authState.user) return;
  const res = await fetch(`${API_BASE}/users/refresh`, {
    method: "POST",
    body: JSON.stringify({
      user_id: authState.user.id,
      refresh_token: authState.user.refresh_token,
    }),
  });
  const json: unknown = await res.json();
  if (!res.ok) {
    authState.logout();
    return new Error("failed to refresh token");
  }
  const parsedJson = Schemas.UserResponseSchema.safeParse(json);
  if (!parsedJson.success) {
    return parsedJson.error;
  }
  authState.setUser(parsedJson.data.data);
}

async function zodFetch<T>(
  url: string,
  schema: z.Schema<T>,
  options: RequestInit | undefined,
): Promise<T | Error> {
  let res = await fetch(url, options);
  const json = await res.json();

  if (res.status === http.StatusUnauthorized) {
    const err = await refreshUserToken();
    if (err instanceof Error) {
      authState.logout();
      return err;
    }

    res = await fetch(url, {
      ...options,
      headers: {
        ...options?.headers,
        Authorization: `Bearer ${authState.user?.access_token}`,
      },
    });
  }
  if (!res.ok) {
    return json.message;
  }
  const parsedResponse = schema.safeParse(json);
  if (!parsedResponse.success) {
    return parsedResponse.error;
  }
  return parsedResponse.data;
}

export const ApiClient = {
  refresh: refreshUserToken,
  zodFetch,
} as const;
