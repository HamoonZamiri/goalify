import { reactive, ref } from "vue";
import type { Goal, GoalCategory } from "@/utils/schemas";

const categoryState = reactive<{ categories: GoalCategory[] }>({
	categories: [],
});

function useGoals() {
	function setCategories(categories: GoalCategory[]) {
		categoryState.categories = categories;
	}

	function addCategory(category: GoalCategory) {
		categoryState.categories.push(category);
	}

	function addGoal(categoryId: string, goal: Goal) {
		const category = categoryState.categories.find((c) => c.id === categoryId);
		if (category) {
			category.goals.push(goal);
		}
	}

	function deleteGoal(categoryId: string, goalId: string) {
		const category = categoryState.categories.find((c) => c.id === categoryId);
		if (category) {
			category.goals = category.goals.filter((g) => g.id !== goalId);
		}
	}
	function deleteCategory(categoryId: string) {
		categoryState.categories = categoryState.categories.filter(
			(c) => c.id !== categoryId,
		);
	}

	return {
		categoryState,
		addCategory,
		addGoal,
		deleteGoal,
		deleteCategory,
		setCategories,
	};
}

export default useGoals;
