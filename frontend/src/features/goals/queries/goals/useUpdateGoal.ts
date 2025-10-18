import { useMutation, useQueryClient } from "@tanstack/vue-query";
import { categoryKeys } from "../queryKeys";
import {
	getAccessToken,
	createAuthHeaders,
	zodFetch,
} from "@/shared/api/client";
import { isErrorResponse, type ErrorResponse } from "@/shared/schemas";
import { API_BASE, http } from "@/utils/constants";
import { GoalSchema, type Goal, type UpdateGoalFormData } from "../../schemas";
import { toast } from "vue3-toastify";

async function updateGoalQueryDataFn(
	goalId: string,
	data: UpdateGoalFormData,
	token: string,
): Promise<Goal | ErrorResponse> {
	return zodFetch(`${API_BASE}/goals/${goalId}`, GoalSchema, {
		method: http.MethodPut,
		headers: createAuthHeaders(token),
		body: JSON.stringify(data),
	});
}

/**
 * Mutation hook to update a goal
 */
export function useUpdateGoal() {
	const queryClient = useQueryClient();
	const token = getAccessToken();

	return useMutation({
		mutationFn: async ({
			goalId,
			data,
		}: {
			goalId: string;
			data: UpdateGoalFormData;
		}) => {
			if (!token) throw new Error("No access token");
			const result = await updateGoalQueryDataFn(goalId, data, token);

			if (isErrorResponse(result)) {
				throw new Error(result.message);
			}

			return result;
		},
		onSuccess: (_) => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
		},
		onError: (error: Error) => {
			toast.error(`Failed to update goal: ${error.message}`);
		},
	});
}
