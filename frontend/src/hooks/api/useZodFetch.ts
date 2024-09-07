import { API_BASE, http } from "@/utils/constants";
import useAuth from "../auth/useAuth";
import router from "@/router";
import { Schemas, type ErrorResponse } from "@/utils/schemas";
import type { z } from "zod";

function useZodFetch() {
  const { setUser, getUser, logout, authState } = useAuth();

  async function refreshUserToken(): Promise<void> {
    if (!getUser) return;
    const res = await fetch(`${API_BASE}/users/refresh`, {
      method: "POST",
      body: JSON.stringify({
        user_id: getUser()?.id,
        refresh_token: authState.value?.refresh_token,
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
    const json = await res.json();

    if (res.status === http.StatusUnauthorized) {
      await refreshUserToken();
      res = await fetch(url, {
        ...options,
        headers: {
          ...options?.headers,
          Authorization: `Bearer ${authState.value?.access_token}`,
        },
      });
    }
    if (!res.ok) {
      const error = json as ErrorResponse;
      error.statusCode = res.status;
      return error;
    }
    const parsedResponse = schema.safeParse(json);
    if (!parsedResponse.success) {
      throw parsedResponse.error;
    }
    return parsedResponse.data;
  }

  return { zodFetch };
}

export default useZodFetch;
