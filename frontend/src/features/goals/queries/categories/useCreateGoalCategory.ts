import { useMutation, useQueryClient } from "@tanstack/vue-query";
import { categoryKeys } from "../queryKeys";
import {
	getAccessToken,
	createAuthHeaders,
	zodFetch,
} from "@/shared/api/client";
import { isErrorResponse, type ErrorResponse } from "@/shared/schemas";
import { API_BASE, http } from "@/utils/constants";
import {
	GoalCategorySchema,
	type GoalCategory,
	type CreateGoalCategoryFormData,
} from "../../schemas";
import { toast } from "vue3-toastify";

async function createGoalCategoryQueryDataFn(
	data: CreateGoalCategoryFormData,
	token: string,
): Promise<GoalCategory | ErrorResponse> {
	return zodFetch(`${API_BASE}/goals/categories`, GoalCategorySchema, {
		method: http.MethodPost,
		headers: createAuthHeaders(token),
		body: JSON.stringify(data),
	});
}

/**
 * Mutation hook to create a new goal category
 */
export function useCreateGoalCategory() {
	const queryClient = useQueryClient();
	const token = getAccessToken();

	return useMutation({
		mutationFn: async (data: CreateGoalCategoryFormData) => {
			if (!token) throw new Error("No access token");
			const result = await createGoalCategoryQueryDataFn(data, token);

			if (isErrorResponse(result)) {
				throw new Error(result.message);
			}

			return result;
		},
		onSuccess: (data) => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
			toast.success(`Successfully created category: ${data.title}`);
		},
		onError: (error: Error) => {
			toast.error(`Failed to create category: ${error.message}`);
		},
	});
}
