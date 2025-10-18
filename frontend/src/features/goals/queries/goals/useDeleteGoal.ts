import { useMutation, useQueryClient } from "@tanstack/vue-query";
import { categoryKeys } from "../queryKeys";
import { getAccessToken, createAuthHeaders } from "@/shared/api/client";
import { API_BASE, http } from "@/utils/constants";
import { toast } from "vue3-toastify";

async function deleteGoalQueryDataFn(
	goalId: string,
	token: string,
): Promise<void> {
	const res = await fetch(`${API_BASE}/goals/${goalId}`, {
		method: http.MethodDelete,
		headers: createAuthHeaders(token),
	});

	if (!res.ok) {
		throw new Error("Failed to delete goal");
	}
}

/**
 * Mutation hook to delete a goal
 */
export function useDeleteGoal() {
	const queryClient = useQueryClient();
	const token = getAccessToken();

	return useMutation({
		mutationFn: async (goalId: string) => {
			if (!token) throw new Error("No access token");
			await deleteGoalQueryDataFn(goalId, token);
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
			toast.success("Successfully deleted goal");
		},
		onError: (error: Error) => {
			toast.error(`Failed to delete goal: ${error.message}`);
		},
	});
}
