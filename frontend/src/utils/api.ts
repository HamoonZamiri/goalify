import authState from "@/state/auth";
import { API_BASE, http } from "./constants";
import { Schemas, type Goal } from "./schemas";
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
): Promise<T | string> {
  let res = await fetch(url, options);
  const json = await res.json();

  if (res.status === http.StatusUnauthorized) {
    const err = await refreshUserToken();
    if (err instanceof Error) {
      authState.logout();
      throw err;
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
    return json.message as string;
  }
  const parsedResponse = schema.safeParse(json);
  if (!parsedResponse.success) {
    throw parsedResponse.error;
  }
  return parsedResponse.data;
}

async function createGoalCategory(
  title: string,
  xp_per_goal: number,
): Promise<string | z.infer<typeof Schemas.GoalCategoryResponseSchema>> {
  try {
    const res = await zodFetch(
      `${API_BASE}/goals/categories`,
      Schemas.GoalCategoryResponseSchema,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authState.user?.access_token}`,
        },
        body: JSON.stringify({
          title,
          xp_per_goal,
        }),
      },
    );
    return res;
  } catch (err) {
    console.error(err);
    return "Failed to create goal category.";
  }
}

async function getUserGoalCategories(): Promise<
  string | z.infer<typeof Schemas.GoalCategoryResponseArraySchema>
> {
  try {
    const res = await zodFetch(
      `${API_BASE}/goals/categories`,
      Schemas.GoalCategoryResponseArraySchema,
      { headers: { Authorization: `Bearer ${authState.user?.access_token}` } },
    );
    return res;
  } catch (err) {
    console.error(err);
    return "Failed to get goal categories.";
  }
}

async function createGoal(
  title: string,
  description: string,
  categoryId: string,
) {
  try {
    const res = await zodFetch(
      `${API_BASE}/goals/create`,
      Schemas.GoalResponseSchema,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authState.user?.access_token}`,
        },
        body: JSON.stringify({
          title,
          description,
          category_id: categoryId,
        }),
      },
    );
    return res;
  } catch (err) {
    console.error(err);
    return "Failed to create goal.";
  }
}

export const ApiClient = {
  refresh: refreshUserToken,
  zodFetch,
  createGoalCategory,
  getUserGoalCategories,
  createGoal,
} as const;
