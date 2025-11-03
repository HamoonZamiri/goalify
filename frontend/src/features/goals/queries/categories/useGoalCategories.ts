import { useQuery } from "@tanstack/vue-query";
import { zodFetch } from "@/shared/api/client";
import { isErrorResponse } from "@/shared/schemas";
import { API_BASE } from "@/utils/constants";
import {
	type GoalCategory,
	GoalCategoryResponseArraySchema,
} from "../../schemas";
import { categoryKeys } from "../queryKeys";

async function goalCategoriesQueryDataFn(): Promise<GoalCategory[]> {
	const result = await zodFetch(
		`${API_BASE}/goals/categories`,
		GoalCategoryResponseArraySchema,
	);

	if (isErrorResponse(result)) {
		throw new Error(result.message);
	}

	return result.data;
}

/**
 * Query hook to fetch all goal categories
 */
export function useGoalCategories() {
	return useQuery({
		queryKey: categoryKeys.lists(),
		queryFn: goalCategoriesQueryDataFn,
	});
}
