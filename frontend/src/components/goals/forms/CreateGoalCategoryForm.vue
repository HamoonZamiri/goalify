<script setup lang="ts">
import useApi from "@/hooks/api/useApi";
import useGoals from "@/hooks/goals/useGoals";
import type { ErrorResponse } from "@/utils/schemas";
import { ref } from "vue";

type CreateCategoryForm = {
  title: string;
  xp_per_goal: number;
};
const formData = ref<CreateCategoryForm>({
  title: "",
  xp_per_goal: 0,
});

const error = ref<ErrorResponse>();

const { addCategory } = useGoals();
const { createGoalCategory, isError } = useApi();
const emit = defineEmits(["submit", "close"]);

async function submit() {
  emit("submit", { ...formData.value });
  const res = await createGoalCategory(
    formData.value.title,
    formData.value.xp_per_goal,
  );
  if (isError(res)) {
    error.value = res;
    return;
  }

  formData.value.title = "";
  formData.value.xp_per_goal = 1;
  error.value = undefined;

  // dispatch an event to update the categories
  addCategory(res);
  emit("close");
}
</script>
<template>
  <form
    @submit.prevent="submit"
    class="rounded-lg border bg-gray-800 p-10 w-[95vw] sm:w-[40vw] grid grid-cols-1 gap-4 hover:cursor-default"
  >
    <p class="flex justify-center text-xl text-gray-200">
      Create a New Goal/Task Category
    </p>
    <div class="">
      <label class="text-gray-200">Title:</label>
      <input
        id="title"
        type="text"
        v-model="formData.title"
        class="block h-10 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6 focus:outline-none"
      />
      <p class="text-red-400" v-if="error?.errors?.title">
        {{ error.errors.title }}
      </p>
    </div>
    <div>
      <label class="text-gray-200">XP/goal:</label>
      <input
        id="xp-per-goal"
        v-model="formData.xp_per_goal"
        type="number"
        min="1"
        max="100"
        class="block h-10 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6 focus:outline-none"
      />
      <p class="text-red-400" v-if="error?.errors?.xp_per_goal">
        {{ error.errors.xp_per_goal }}
      </p>
    </div>
    <button
      type="submit"
      class="bg-blue-400 mt-10 rounded-lg text-gray-300 hover:bg-blue-500 h-10 py-1.5"
    >
      Add Category
    </button>
  </form>
</template>
