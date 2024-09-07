import useZodFetch from "@/hooks/api/useZodFetch";
import { API_BASE, http } from "@/utils/constants";
import {
  Schemas,
  type ErrorResponse,
  type Goal,
  type GoalCategory,
} from "@/utils/schemas";
import useAuth from "../auth/useAuth";

type ServerResponse<T> = {
  message: string;
  data: T;
};
function useApi() {
  const { zodFetch } = useZodFetch();
  const { authState } = useAuth();

  function isError(
    res: ServerResponse<any> | ErrorResponse,
  ): res is ErrorResponse {
    const casted = res as ErrorResponse;
    return (
      casted.errors !== undefined ||
      (casted.statusCode !== undefined && casted.statusCode >= 400)
    );
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
          Authorization: `Bearer ${authState.value?.access_token}`,
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
      {
        headers: { Authorization: `Bearer ${authState.value?.access_token}` },
      },
    );
    return res;
  }

  async function createGoal(
    title: string,
    description: string,
    categoryId: string,
  ): Promise<ErrorResponse | ServerResponse<Goal>> {
    const res = await zodFetch(
      `${API_BASE}/goals`,
      Schemas.GoalResponseSchema,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authState.value?.access_token}`,
        },
        body: JSON.stringify({
          title,
          description,
          category_id: categoryId,
        }),
      },
    );
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
          Authorization: `Bearer ${authState.value?.access_token}`,
        },
        body: JSON.stringify(updates),
      },
    );
    return res;
  }

  async function deleteGoal(goalId: string): Promise<void> {
    const res = await fetch(`${API_BASE}/goals/${goalId}`, {
      method: http.MethodDelete,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authState.value?.access_token}`,
      },
    });
    if (!res.ok) {
      throw new Error("Failed to delete goal");
    }
  }

  async function deleteCategory(categoryId: string): Promise<void> {
    const res = await fetch(`${API_BASE}/goals/categories/${categoryId}`, {
      method: http.MethodDelete,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authState.value?.access_token}`,
      },
    });
    if (!res.ok) {
      throw new Error("Failed to delete category");
    }
  }

  async function updateCategory(
    categoryId: string,
    updates: Partial<GoalCategory>,
  ) {
    const res = await fetch(`${API_BASE}/goals/categories/${categoryId}`, {
      method: http.MethodPut,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authState.value?.access_token}`,
      },
      body: JSON.stringify(updates),
    });
    if (!res.ok) {
      throw new Error("Failed to update category");
    }
  }

  return {
    isError,
    createGoalCategory,
    getUserGoalCategories,
    createGoal,
    updateGoal,
    deleteGoal,
    deleteCategory,
    updateCategory,
  };
}

export default useApi;
