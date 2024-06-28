<script setup lang="ts">
import goalState from "@/state/goals";
import { ApiClient } from "@/utils/api";
import { ref } from "vue";

type CreateCategoryForm = {
  title: string;
  xp_per_goal: number;
};
const formData = ref<CreateCategoryForm>({
  title: "",
  xp_per_goal: 0,
});

const error = ref<string | null>(null);

const props = defineProps<{
  isOpen: boolean;
  setIsOpen: (value: boolean) => void;
}>();

async function handleSubmit(e: MouseEvent) {
  e.preventDefault();
  props.setIsOpen(false);
  const res = await ApiClient.createGoalCategory(
    formData.value.title,
    formData.value.xp_per_goal,
  );
  if (typeof res === "string") {
    error.value = res;
    return;
  }

  formData.value.title = "";
  formData.value.xp_per_goal = 0;

  // dispatch an event to update the categories
  goalState.addCategory(res.data);
  props.setIsOpen(false);
}
</script>
<template>
  <form
    class="rounded-lg border bg-white p-3 w-[400px] grid grid-cols-1 gap-8 hover:cursor-default"
  >
    <p class="font-semibold">Create a New Goal/Task Category</p>
    <div class="grid grid-cols-1 gap-4">
      <label>Title:</label>
      <input
        type="text"
        v-model="formData.title"
        class="block w-full rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-200 sm:text-sm sm:leading-6"
      />
    </div>
    <div>
      <label>XP/goal:</label>
      <input
        v-model="formData.xp_per_goal"
        type="number"
        min="1"
        max="100"
        class="block w-full rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-200 sm:text-sm sm:leading-6"
      />
    </div>
    <button
      @click="handleSubmit"
      class="block w-full hover:bg-blue-200 bg-blue-100 rounded-lg h-10"
    >
      Add Category
    </button>
  </form>
</template>
