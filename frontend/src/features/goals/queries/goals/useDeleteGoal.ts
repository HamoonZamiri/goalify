import {
	type UseMutationOptions,
	useMutation,
	useQueryClient,
} from "@tanstack/vue-query";
import { API_BASE, http } from "@/utils/constants";
import { categoryKeys } from "../queryKeys";

async function deleteGoalQueryDataFn(goalId: string): Promise<void> {
	const res = await fetch(`${API_BASE}/goals/${goalId}`, {
		method: http.MethodDelete,
		headers: {
			"Content-Type": "application/json",
		},
	});

	if (!res.ok) {
		const json = await res.json();
		throw new Error(json.message || "Failed to delete goal");
	}
}

/**
 * Mutation hook to delete a goal
 */
export function useDeleteGoal(
	options?: UseMutationOptions<void, Error, string>,
) {
	const queryClient = useQueryClient();

	return useMutation({
		...options,
		mutationFn: deleteGoalQueryDataFn,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
		},
	});
}
