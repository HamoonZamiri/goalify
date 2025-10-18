import {
	type UseMutationOptions,
	useMutation,
	useQueryClient,
} from "@tanstack/vue-query";
import { zodFetch } from "@/shared/api/client";
import { isErrorResponse } from "@/shared/schemas";
import { API_BASE, http } from "@/utils/constants";
import { type Goal, GoalSchema, type UpdateGoalFormData } from "../../schemas";
import { categoryKeys } from "../queryKeys";

type UpdateGoalVariables = {
	goalId: string;
	data: UpdateGoalFormData;
};

async function updateGoalQueryDataFn(
	goalId: string,
	data: UpdateGoalFormData,
): Promise<Goal> {
	const result = await zodFetch(`${API_BASE}/goals/${goalId}`, GoalSchema, {
		method: http.MethodPut,
		body: JSON.stringify(data),
	});

	if (isErrorResponse(result)) {
		throw new Error(result.message);
	}

	return result;
}

/**
 * Mutation hook to update a goal
 */
export function useUpdateGoal(
	options?: UseMutationOptions<Goal, Error, UpdateGoalVariables>,
) {
	const queryClient = useQueryClient();

	return useMutation({
		...options,
		mutationFn: ({ goalId, data }: UpdateGoalVariables) =>
			updateGoalQueryDataFn(goalId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
		},
	});
}
