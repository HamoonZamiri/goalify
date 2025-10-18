import {
	type UseMutationOptions,
	useMutation,
	useQueryClient,
} from "@tanstack/vue-query";
import { zodFetch } from "@/shared/api/client";
import { isErrorResponse } from "@/shared/schemas";
import { API_BASE, http } from "@/utils/constants";
import {
	type CreateGoalCategoryFormData,
	type GoalCategory,
	GoalCategorySchema,
} from "../../schemas";
import { categoryKeys } from "../queryKeys";

async function createGoalCategoryQueryDataFn(
	data: CreateGoalCategoryFormData,
): Promise<GoalCategory> {
	const result = await zodFetch(
		`${API_BASE}/goals/categories`,
		GoalCategorySchema,
		{
			method: http.MethodPost,
			body: JSON.stringify(data),
		},
	);

	if (isErrorResponse(result)) {
		throw new Error(result.message);
	}

	return result;
}

/**
 * Mutation hook to create a new goal category
 */
export function useCreateGoalCategory(
	options?: UseMutationOptions<GoalCategory, Error, CreateGoalCategoryFormData>,
) {
	const queryClient = useQueryClient();

	return useMutation({
		...options,
		mutationFn: createGoalCategoryQueryDataFn,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
		},
	});
}
