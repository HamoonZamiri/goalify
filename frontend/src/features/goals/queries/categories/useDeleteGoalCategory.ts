import {
	type UseMutationOptions,
	useMutation,
	useQueryClient,
} from "@tanstack/vue-query";
import { z } from "zod";
import { API_BASE, http } from "@/utils/constants";
import { categoryKeys } from "../queryKeys";

const DeleteResponseSchema = z.object({ message: z.string() });

async function deleteGoalCategoryQueryDataFn(
	categoryId: string,
): Promise<void> {
	const res = await fetch(`${API_BASE}/goals/categories/${categoryId}`, {
		method: http.MethodDelete,
		headers: {
			"Content-Type": "application/json",
		},
	});

	if (!res.ok) {
		const json = await res.json();
		throw new Error(json.message || "Failed to delete category");
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
