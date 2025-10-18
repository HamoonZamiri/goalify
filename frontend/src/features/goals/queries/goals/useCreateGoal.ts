import {
	type UseMutationOptions,
	useMutation,
	useQueryClient,
} from "@tanstack/vue-query";
import { zodFetch } from "@/shared/api/client";
import { isErrorResponse } from "@/shared/schemas";
import { API_BASE, http } from "@/utils/constants";
import { type CreateGoalFormData, type Goal, GoalSchema } from "../../schemas";
import { categoryKeys } from "../queryKeys";

async function createGoalQueryDataFn(data: CreateGoalFormData): Promise<Goal> {
	const result = await zodFetch(`${API_BASE}/goals`, GoalSchema, {
		method: http.MethodPost,
		body: JSON.stringify(data),
	});

	if (isErrorResponse(result)) {
		throw new Error(result.message);
	}

	return result;
}

/**
 * Mutation hook to create a new goal
 */
export function useCreateGoal(
	options?: UseMutationOptions<Goal, Error, CreateGoalFormData>,
) {
	const queryClient = useQueryClient();

	return useMutation({
		...options,
		mutationFn: createGoalQueryDataFn,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
		},
	});
}
