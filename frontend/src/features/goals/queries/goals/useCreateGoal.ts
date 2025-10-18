import { useMutation, useQueryClient } from "@tanstack/vue-query";
import { categoryKeys } from "../queryKeys";
import {
	getAccessToken,
	createAuthHeaders,
	zodFetch,
} from "@/shared/api/client";
import { isErrorResponse, type ErrorResponse } from "@/shared/schemas";
import { API_BASE, http } from "@/utils/constants";
import { GoalSchema, type Goal, type CreateGoalFormData } from "../../schemas";
import { toast } from "vue3-toastify";

async function createGoalQueryDataFn(
	data: CreateGoalFormData,
	token: string,
): Promise<Goal | ErrorResponse> {
	return zodFetch(`${API_BASE}/goals`, GoalSchema, {
		method: http.MethodPost,
		headers: createAuthHeaders(token),
		body: JSON.stringify(data),
	});
}

/**
 * Mutation hook to create a new goal
 */
export function useCreateGoal() {
	const queryClient = useQueryClient();
	const token = getAccessToken();

	return useMutation({
		mutationFn: async (data: CreateGoalFormData) => {
			if (!token) throw new Error("No access token");
			const result = await createGoalQueryDataFn(data, token);

			if (isErrorResponse(result)) {
				throw new Error(result.message);
			}

			return result;
		},
		onSuccess: (data) => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
			toast.success(`Successfully created goal: ${data.title}`);
		},
		onError: (error: Error) => {
			toast.error(`Failed to create goal: ${error.message}`);
		},
	});
}
