<script setup lang="ts">
import type { GoalCategory } from "@/utils/schemas";
import GoalCard from "./GoalCard.vue";
import ModalForm from "@/components/ModalForm.vue";
import CreateGoalForm from "@/components/goals/forms/CreateGoalForm.vue";
import CreateGoalButton from "@/components/goals/buttons/CreateGoalButton.vue";
import { Menu, MenuButton, MenuItems, MenuItem } from "@headlessui/vue";
import { reactive, ref, watch } from "vue";
import useGoals from "@/hooks/goals/useGoals";
import useApi from "@/hooks/api/useApi";
import { toast } from "vue3-toastify";
const props = defineProps<{
  goalCategory: GoalCategory;
}>();

const XP_PER_GOAL_MAX = 100;
const XP_PER_GOAL_MIN = 1;

const isCreateGoalDialogOpen = ref(false);
const updates = reactive({
  title: props.goalCategory.title,
  xp_per_goal: props.goalCategory.xp_per_goal,
});

const { deleteCategory } = useGoals();
const { deleteCategory: apiDeleteCategory, updateCategory: apiUpdateCategory } =
  useApi();

async function handleDeleteCategory(e: MouseEvent) {
  e.preventDefault();
  await apiDeleteCategory(props.goalCategory.id);

  // remove category from state
  deleteCategory(props.goalCategory.id);

  toast.success(`Successfully deleted category: ${props.goalCategory.title}`);
}

async function handleNumericInput(payload: Event) {
  const value = (payload.target as HTMLInputElement).value;
  const parsedVal = Number.parseInt(value);
  if (parsedVal > XP_PER_GOAL_MAX) {
    updates.xp_per_goal = XP_PER_GOAL_MAX;
  } else {
    updates.xp_per_goal = parsedVal;
  }
}

watch(updates, async (category) => {
  props.goalCategory.title = category.title;
  props.goalCategory.xp_per_goal = category.xp_per_goal;

  if (
    category.title === "" ||
    category.xp_per_goal < XP_PER_GOAL_MIN ||
    category.xp_per_goal > XP_PER_GOAL_MAX ||
    Number.isNaN(category.xp_per_goal)
  )
    return;

  await apiUpdateCategory(props.goalCategory.id, {
    title: category.title,
    xp_per_goal: category.xp_per_goal,
  });
});
</script>
<template>
  <div class="flex flex-col">
    <header class="flex justify-between">
      <div class="flex flex-col">
        <input
          v-model="updates.title"
          class="w-auto text-gray-200 bg-gray-900 text-xl focus:outline-none"
        />
        <div class="flex gap-1">
          <span class="text-xs text-gray-300">Earn</span>
          <input
            type="number"
            min="1"
            max="100"
            @input="handleNumericInput"
            class="num-input focus:outline-none text-xs bg-gray-900 font-semibold text-green-500"
            :class="{
              'w-2':
                props.goalCategory.xp_per_goal < 10 ||
                isNaN(props.goalCategory.xp_per_goal),
              'w-4':
                props.goalCategory.xp_per_goal >= 10 &&
                props.goalCategory.xp_per_goal < 100,
              'w-6': props.goalCategory.xp_per_goal == 100,
            }"
            v-model="updates.xp_per_goal"
          />
          <span class="text-xs text-gray-300">xp/goal</span>
        </div>
      </div>
      <div class="flex">
        <CreateGoalButton
          class="hover:cursor-pointer"
          @click="isCreateGoalDialogOpen = true"
        />
        <ModalForm v-model="isCreateGoalDialogOpen">
          <CreateGoalForm
            :category-id="props.goalCategory.id"
            @close="isCreateGoalDialogOpen = false"
          />
        </ModalForm>
        <Menu as="div" class="relative inline-block">
          <MenuButton class="hover:cursor-pointer">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              class="size-6 stroke-gray-300"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M8.625 12a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H8.25m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H12m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0h-.375M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
              />
            </svg>
          </MenuButton>
          <MenuItems
            class="absolute flex flex-col items-start w-56 bg-gray-500 p-1 rounded-md justify-self-start right-0"
          >
            <MenuItem
              as="button"
              class="w-full px-2 flex justify-start gap-x-2 bg-gray-600 hover:cursor-pointer hover:bg-gray-400 text-gray-700"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                class="size-5 stroke-gray-300"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L6.832 19.82a4.5 4.5 0 0 1-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 0 1 1.13-1.897L16.863 4.487Zm0 0L19.5 7.125"
                />
              </svg>
              <span class="text-sm text-gray-300">Edit</span>
            </MenuItem>

            <MenuItem
              as="button"
              class="w-full px-2 hover:cursor-pointer hover:bg-gray-400 bg-gray-600 text-gray-700 flex gap-x-2"
              @click="handleDeleteCategory"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                class="size-5 stroke-gray-300"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
                />
              </svg>
              <span class="text-sm text-gray-300">Delete</span>
            </MenuItem>
          </MenuItems>
        </Menu>
      </div>
    </header>
    <div class="w-full flex flex-col gap-2" v-for="goal in goalCategory.goals">
      <GoalCard :goal="goal" />
    </div>
  </div>
</template>
<style scoped>
/* Hide default increment and decrement arrows */
.num-input::-webkit-outer-spin-button,
.num-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.num-input {
  -moz-appearance: textfield; /* Hides the arrows in Firefox */
  appearance: textfield; /* Hides the arrows in other browsers */
}
</style>
