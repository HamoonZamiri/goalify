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
	type UpdateGoalCategoryFormData,
} from "../../schemas";
import { toast } from "vue3-toastify";

async function updateGoalCategoryQueryDataFn(
	categoryId: string,
	data: UpdateGoalCategoryFormData,
	token: string,
): Promise<GoalCategory | ErrorResponse> {
	return zodFetch(
		`${API_BASE}/goals/categories/${categoryId}`,
		GoalCategorySchema,
		{
			method: http.MethodPut,
			headers: createAuthHeaders(token),
			body: JSON.stringify(data),
		},
	);
}

/**
 * Mutation hook to update a goal category
 */
export function useUpdateGoalCategory() {
	const queryClient = useQueryClient();
	const token = getAccessToken();

	return useMutation({
		mutationFn: async ({
			categoryId,
			data,
		}: {
			categoryId: string;
			data: UpdateGoalCategoryFormData;
		}) => {
			if (!token) throw new Error("No access token");
			const result = await updateGoalCategoryQueryDataFn(
				categoryId,
				data,
				token,
			);

			if (isErrorResponse(result)) {
				throw new Error(result.message);
			}

			return result;
		},
		onSuccess: (data) => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
			toast.success(`Successfully updated category: ${data.title}`);
		},
		onError: (error: Error) => {
			toast.error(`Failed to update category: ${error.message}`);
		},
	});
}
