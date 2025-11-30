import {
	type UseMutationOptions,
	useMutation,
	useQueryClient,
} from "@tanstack/vue-query";
import { toast } from "vue3-toastify";
import { z } from "zod";
import { zodFetch } from "@/shared/api/client";
import { isErrorResponse } from "@/shared/schemas";
import { API_BASE, http } from "@/utils/constants";
import { ResetGoalCategoryParams } from "../../schemas";
import { categoryKeys } from "../queryKeys";

async function resetGoalCategoryQueryDataFn(
	data: ResetGoalCategoryParams,
): Promise<void> {
	const result = await zodFetch(
		`${API_BASE}/goals/categories/${data.category_id}/reset`,
		z.object({}), // Empty schema - zodFetch handles 204 No Content
		{
			method: http.MethodPost,
		},
	);

	if (isErrorResponse(result)) {
		throw new Error(result.message);
	}
}

/**
 * Mutation hook to create a new goal category
 */
export function useResetGoalCategory(
	options?: UseMutationOptions<void, Error, ResetGoalCategoryParams>,
) {
	const queryClient = useQueryClient();

	return useMutation({
		...options,
		mutationFn: resetGoalCategoryQueryDataFn,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
		},
	});
}
