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

const goalState = reactive<{
  categories: GoalCategory[];
  addCategory: (category: GoalCategory) => void;
  addGoal: (categoryId: string, goal: Goal) => void;
}>({
  categories: [],
  addCategory,
  addGoal,
});

export default goalState;
