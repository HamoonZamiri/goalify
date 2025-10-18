import { useQuery } from "@tanstack/vue-query";
import { categoryKeys } from "../queryKeys";
import {
	getAccessToken,
	createAuthHeaders,
	zodFetch,
} from "@/shared/api/client";
import { isErrorResponse, type ErrorResponse } from "@/shared/schemas";
import { API_BASE } from "@/utils/constants";
import {
	GoalCategoryResponseArraySchema,
	type GoalCategory,
} from "../../schemas";

async function goalCategoriesQueryDataFn(
	token: string,
): Promise<{ data: GoalCategory[] } | ErrorResponse> {
	return zodFetch(
		`${API_BASE}/goals/categories`,
		GoalCategoryResponseArraySchema,
		{
			headers: createAuthHeaders(token),
		},
	);
}

/**
 * Query hook to fetch all goal categories
 */
export function useGoalCategories() {
	const token = getAccessToken();

	return useQuery({
		queryKey: categoryKeys.lists(),
		queryFn: async () => {
			if (!token) throw new Error("No access token");
			const result = await goalCategoriesQueryDataFn(token);

			if (isErrorResponse(result)) {
				throw new Error(result.message);
			}

			return result.data;
		},
		enabled: !!token,
	});
}
