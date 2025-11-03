import {
  type UseMutationOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/vue-query";
import { API_BASE, http } from "@/utils/constants";
import { categoryKeys } from "../queryKeys";
import { z } from "zod";
import { zodFetch } from "@/shared/api";
import { isErrorResponse } from "@/shared/schemas";

async function deleteGoalQueryDataFn(goalId: string): Promise<void> {
  const res = await zodFetch(`${API_BASE}/goals/${goalId}`, z.object({}), {
    method: http.MethodDelete,
  });

  if (isErrorResponse(res)) {
    throw new Error(res.message);
  }
}

/**
 * Mutation hook to delete a goal
 */
export function useDeleteGoal(
  options?: UseMutationOptions<void, Error, string>,
) {
  const queryClient = useQueryClient();

  return useMutation({
    ...options,
    mutationFn: deleteGoalQueryDataFn,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: categoryKeys.all });
    },
  });
}
