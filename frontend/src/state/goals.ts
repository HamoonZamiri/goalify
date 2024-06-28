import { reactive } from "vue";
import { type GoalCategory } from "@/utils/schemas";

function addCategory(category: GoalCategory) {
  goalState.categories.push(category);
}

const goalState = reactive<{
  categories: GoalCategory[];
  addCategory: (category: GoalCategory) => void;
}>({
  categories: [],
  addCategory,
});

export default goalState;
