import { API_BASE, http } from "@/utils/constants";
import useAuth from "../auth/useAuth";
import router from "@/router";
import { Schemas, type ErrorResponse, type User } from "@/utils/schemas";
import type { z } from "zod";

let refreshPromise: Promise<User | undefined> | undefined = undefined;

function useZodFetch() {
  const { setUser, getUser, logout } = useAuth();

  async function refreshUserToken(): Promise<User | undefined> {
    if (refreshPromise) {
      return refreshPromise;
    }

    refreshPromise = (async () => {
      try {
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

  async function zodFetch<T>(
    url: string,
    schema: z.Schema<T>,
    options: RequestInit | undefined,
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
          Authorization: `Bearer ${refreshedUser.access_token}`,
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

  return { zodFetch };
}

export default useZodFetch;
