import { API_BASE, http } from "@/utils/constants";
import useAuth from "../auth/useAuth";
import router from "@/router";
import { Schemas, type ErrorResponse } from "@/utils/schemas";
import type { z } from "zod";

function useZodFetch() {
  const { setUser, getUser, logout, authState } = useAuth();

  async function refreshUserToken(): Promise<void> {
    const user = getUser();
    if (!user) return;
    const res = await fetch(`${API_BASE}/users/refresh`, {
      method: "POST",
      body: JSON.stringify({
        user_id: user.id,
        refresh_token: user.refresh_token,
      }),
    });
    const json: unknown = await res.json();
    if (!res.ok) {
      logout();
      router.push({ name: "Login" });
      return;
    }
    const parsedJson = Schemas.UserResponseSchema.parse(json);
    setUser(parsedJson.data);
  }

  async function zodFetch<T>(
    url: string,
    schema: z.Schema<T>,
    options: RequestInit | undefined,
  ): Promise<T | ErrorResponse> {
    let res = await fetch(url, options);
    let json = await res.json();

    if (res.status === http.StatusUnauthorized) {
      await refreshUserToken();
      const user = getUser();
      if (!user) return json;
      res = await fetch(url, {
        ...options,
        headers: {
          ...options?.headers,
          Authorization: `Bearer ${user.access_token}`,
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
