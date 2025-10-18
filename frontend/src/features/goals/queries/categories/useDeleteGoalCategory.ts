import { useMutation, useQueryClient } from "@tanstack/vue-query";
import { categoryKeys } from "../queryKeys";
import { getAccessToken, createAuthHeaders } from "@/shared/api/client";
import { API_BASE, http } from "@/utils/constants";
import { toast } from "vue3-toastify";

async function deleteGoalCategoryQueryDataFn(
	categoryId: string,
	token: string,
): Promise<void> {
	const res = await fetch(`${API_BASE}/goals/categories/${categoryId}`, {
		method: http.MethodDelete,
		headers: createAuthHeaders(token),
	});

	if (!res.ok) {
		throw new Error("Failed to delete category");
	}
}

/**
 * Mutation hook to delete a goal category
 */
export function useDeleteGoalCategory() {
	const queryClient = useQueryClient();
	const token = getAccessToken();

	return useMutation({
		mutationFn: async (categoryId: string) => {
			if (!token) throw new Error("No access token");
			await deleteGoalCategoryQueryDataFn(categoryId, token);
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
			toast.success("Successfully deleted category");
		},
		onError: (error: Error) => {
			toast.error(`Failed to delete category: ${error.message}`);
		},
	});
}
