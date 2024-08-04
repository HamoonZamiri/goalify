import { reactive } from "vue";
import { type Goal, type GoalCategory } from "@/utils/schemas";

function addCategory(category: GoalCategory) {
  goalState.categories.push(category);
}

function addGoal(categoryId: string, goal: Goal) {
  const category = goalState.categories.find((c) => c.id === categoryId);
  if (category) {
    category.goals.push(goal);
  }
}

function deleteGoal(categoryId: string, goalId: string) {
  const category = goalState.categories.find((c) => c.id === categoryId);
  if (category) {
    category.goals = category.goals.filter((g) => g.id !== goalId);
  }
}

function deleteCategory(categoryId: string) {
  goalState.categories = goalState.categories.filter(
    (c) => c.id !== categoryId,
  );
}

const goalState = reactive<{
  categories: GoalCategory[];
  addCategory: (category: GoalCategory) => void;
  addGoal: (categoryId: string, goal: Goal) => void;
  deleteGoal: (categoryId: string, goalId: string) => void;
  deleteCategory: (categoryId: string) => void;
}>({
  categories: [],
  addCategory,
  addGoal,
  deleteGoal,
  deleteCategory,
});

export default goalState;
