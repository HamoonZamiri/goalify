import useZodFetch from "@/hooks/api/useZodFetch";
import { API_BASE, http } from "@/utils/constants";
import {
	type ErrorResponse,
	type Goal,
	type GoalCategory,
	Schemas,
} from "@/utils/schemas";
import useAuth from "../auth/useAuth";

type ServerResponse<T> = {
	data: T;
};
function useApi() {
	const { zodFetch } = useZodFetch();
	const { authState } = useAuth();

	function isError(res: unknown | ErrorResponse): res is ErrorResponse {
		const casted = res as ErrorResponse;
		return (
			casted.errors !== undefined ||
			(casted.statusCode !== undefined && casted.statusCode >= 400)
		);
	}

	async function createGoalCategory(
		title: string,
		xp_per_goal: number,
	): Promise<ErrorResponse | GoalCategory> {
		const res = await zodFetch(
			`${API_BASE}/goals/categories`,
			Schemas.GoalCategorySchema,
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					Authorization: `Bearer ${authState.value?.access_token}`,
				},
				body: JSON.stringify({
					title,
					xp_per_goal,
				}),
			},
		);
		return res;
	}

	async function getUserGoalCategories(): Promise<
		ErrorResponse | ServerResponse<GoalCategory[]>
	> {
		const res = await zodFetch(
			`${API_BASE}/goals/categories`,
			Schemas.GoalCategoryResponseArraySchema,
			{
				headers: { Authorization: `Bearer ${authState.value?.access_token}` },
			},
		);
		return res;
	}

	async function createGoal(
		title: string,
		description: string,
		categoryId: string,
	): Promise<ErrorResponse | Goal> {
		const res = await zodFetch(`${API_BASE}/goals`, Schemas.GoalSchema, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				Authorization: `Bearer ${authState.value?.access_token}`,
			},
			body: JSON.stringify({
				title,
				description,
				category_id: categoryId,
			}),
		});
		return res;
	}

	async function updateGoal(
		goalId: string,
		updates: Partial<Goal>,
	): Promise<ErrorResponse | Goal> {
		const res = await zodFetch(
			`${API_BASE}/goals/${goalId}`,
			Schemas.GoalSchema,
			{
				method: http.MethodPut,
				headers: {
					"Content-Type": "application/json",
					Authorization: `Bearer ${authState.value?.access_token}`,
				},
				body: JSON.stringify(updates),
			},
		);
		return res;
	}

	async function deleteGoal(goalId: string): Promise<void> {
		const res = await fetch(`${API_BASE}/goals/${goalId}`, {
			method: http.MethodDelete,
			headers: {
				"Content-Type": "application/json",
				Authorization: `Bearer ${authState.value?.access_token}`,
			},
		});
		if (!res.ok) {
			throw new Error("Failed to delete goal");
		}
	}

	async function deleteCategory(categoryId: string): Promise<void> {
		const res = await fetch(`${API_BASE}/goals/categories/${categoryId}`, {
			method: http.MethodDelete,
			headers: {
				"Content-Type": "application/json",
				Authorization: `Bearer ${authState.value?.access_token}`,
			},
		});
		if (!res.ok) {
			throw new Error("Failed to delete category");
		}
	}

	async function updateCategory(
		categoryId: string,
		updates: Partial<GoalCategory>,
	) {
		const res = await zodFetch(
			`${API_BASE}/goals/categories/${categoryId}`,
			Schemas.GoalCategorySchema,
			{
				method: http.MethodPut,
				headers: {
					"Content-Type": "application/json",
					Authorization: `Bearer ${authState.value?.access_token}`,
				},
				body: JSON.stringify(updates),
			},
		);
		return res;
	}

	async function getLevel(id: number) {
		const res = await zodFetch(
			`${API_BASE}/levels/${id}`,
			Schemas.LevelSchema,
			{
				headers: { Authorization: `Bearer ${authState.value?.access_token}` },
			},
		);
		return res;
	}

	return {
		isError,
		createGoalCategory,
		getUserGoalCategories,
		createGoal,
		updateGoal,
		deleteGoal,
		deleteCategory,
		updateCategory,
		getLevel,
	};
}

export default useApi;
