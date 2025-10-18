import {
	type UseMutationOptions,
	useMutation,
	useQueryClient,
} from "@tanstack/vue-query";
import { zodFetch } from "@/shared/api/client";
import { isErrorResponse } from "@/shared/schemas";
import { API_BASE, http } from "@/utils/constants";
import {
	type GoalCategory,
	GoalCategorySchema,
	type UpdateGoalCategoryFormData,
} from "../../schemas";
import { categoryKeys } from "../queryKeys";

type UpdateGoalCategoryVariables = {
	categoryId: string;
	data: UpdateGoalCategoryFormData;
};

async function updateGoalCategoryQueryDataFn(
	categoryId: string,
	data: UpdateGoalCategoryFormData,
): Promise<GoalCategory> {
	const result = await zodFetch(
		`${API_BASE}/goals/categories/${categoryId}`,
		GoalCategorySchema,
		{
			method: http.MethodPut,
			body: JSON.stringify(data),
		},
	);

	if (isErrorResponse(result)) {
		throw new Error(result.message);
	}

	return result;
}

/**
 * Mutation hook to update a goal category
 */
export function useUpdateGoalCategory(
	options?: UseMutationOptions<
		GoalCategory,
		Error,
		UpdateGoalCategoryVariables
	>,
) {
	const queryClient = useQueryClient();

	return useMutation({
		...options,
		mutationFn: ({ categoryId, data }: UpdateGoalCategoryVariables) =>
			updateGoalCategoryQueryDataFn(categoryId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
		},
	});
}
