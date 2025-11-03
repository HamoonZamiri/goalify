import {
	type UseMutationOptions,
	useMutation,
	useQueryClient,
} from "@tanstack/vue-query";
import { z } from "zod";
import { zodFetch } from "@/shared/api/client";
import { isErrorResponse } from "@/shared/schemas";
import { API_BASE, http } from "@/utils/constants";
import { categoryKeys } from "../queryKeys";

async function deleteGoalCategoryQueryDataFn(
	categoryId: string,
): Promise<void> {
	const res = await zodFetch(
		`${API_BASE}/goals/categories/${categoryId}`,
		z.object({}),
		{
			method: http.MethodDelete,
		},
	);
	if (isErrorResponse(res)) {
		throw new Error(res.message);
	}
}

/**
 * Mutation hook to delete a goal category
 */
export function useDeleteGoalCategory(
	options?: UseMutationOptions<void, Error, string>,
) {
	const queryClient = useQueryClient();

	return useMutation({
		...options,
		mutationFn: deleteGoalCategoryQueryDataFn,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
		},
	});
}
