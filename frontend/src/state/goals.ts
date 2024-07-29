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

const goalState = reactive<{
  categories: GoalCategory[];
  addCategory: (category: GoalCategory) => void;
  addGoal: (categoryId: string, goal: Goal) => void;
  deleteGoal: (categoryId: string, goalId: string) => void;
}>({
  categories: [],
  addCategory,
  addGoal,
  deleteGoal,
});

export default goalState;
