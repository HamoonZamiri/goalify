import authState from "@/state/auth";
import { API_BASE, http } from "./constants";
import { Schemas, type Goal, type GoalCategory } from "./schemas";
import type { z } from "zod";
import router, { RouteNames } from "@/router";

type ServerResponse<T> = {
  message: string;
  data: T;
};

type ErrorMap = Record<string, string>;
export type ErrorResponse = {
  // we will manually add this field to our errors from the json response
  statusCode?: number;
  // message should always be present
  message: string;
  // in creation requests the server returns an object mapping field names to error messages
  errors?: ErrorMap;
};

function isError(
  res: ServerResponse<any> | ErrorResponse,
): res is ErrorResponse {
  const casted = res as ErrorResponse;
  return (
    casted.errors !== undefined ||
    (casted.statusCode !== undefined && casted.statusCode >= 400)
  );
}

async function refreshUserToken(): Promise<void> {
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
    router.push({ name: "Login" });
    return;
  }
  const parsedJson = Schemas.UserResponseSchema.safeParse(json);
  if (!parsedJson.success) {
    throw parsedJson.error;
  }
  authState.setUser(parsedJson.data.data);
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
        Authorization: `Bearer ${authState.user?.access_token}`,
      },
    });
    router.push(RouteNames.HOME);
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

async function createGoalCategory(
  title: string,
  xp_per_goal: number,
): Promise<ErrorResponse | ServerResponse<GoalCategory>> {
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
}

async function getUserGoalCategories(): Promise<
  ErrorResponse | ServerResponse<GoalCategory[]>
> {
  const res = await zodFetch(
    `${API_BASE}/goals/categories`,
    Schemas.GoalCategoryResponseArraySchema,
    { headers: { Authorization: `Bearer ${authState.user?.access_token}` } },
  );
  return res;
}

async function createGoal(
  title: string,
  description: string,
  categoryId: string,
): Promise<ErrorResponse | ServerResponse<Goal>> {
  const res = await zodFetch(`${API_BASE}/goals`, Schemas.GoalResponseSchema, {
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
  });
  return res;
}

async function updateGoal(
  goalId: string,
  updates: Partial<Goal>,
): Promise<ErrorResponse | ServerResponse<Goal>> {
  const res = await zodFetch(
    `${API_BASE}/goals/${goalId}`,
    Schemas.GoalResponseSchema,
    {
      method: http.MethodPut,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authState.user?.access_token}`,
      },
      body: JSON.stringify(updates),
    },
  );
  return res;
}

// async function deleteGoal(goalId: string)

export const ApiClient = {
  refresh: refreshUserToken,
  zodFetch,
  createGoalCategory,
  getUserGoalCategories,
  createGoal,
  updateGoal,
  isError,
} as const;
