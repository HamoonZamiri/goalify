<script setup lang="ts">
import useApi from "@/hooks/api/useApi";
import useGoals from "@/hooks/goals/useGoals";
import type { ErrorResponse } from "@/utils/schemas";
import { ref } from "vue";

type CreateGoalForm = {
  title: string;
  description: string;
};

const formData = ref<CreateGoalForm>({
  title: "",
  description: "",
});

const error = ref<ErrorResponse | null>(null);
const { addGoal } = useGoals();
const { createGoal, isError } = useApi();

const CreateGoalFormProps = defineProps<{
  props: {
    categoryId: string;
  };
  isOpen: boolean;
  setIsOpen: (value: boolean) => void;
}>();

async function handleSubmit(e: MouseEvent) {
  e.preventDefault();
  const { title, description } = formData.value;
  const res = await createGoal(
    title,
    description,
    CreateGoalFormProps.props.categoryId,
  );
  if (isError(res)) {
    error.value = res;
    return;
  }
  addGoal(CreateGoalFormProps.props.categoryId, res.data);
  formData.value.title = "";
  formData.value.description = "";
  CreateGoalFormProps.setIsOpen(false);
}
</script>

<template>
  <form
    class="rounded-lg border bg-gray-800 p-6 w-[95vw] sm:w-[40vw] grid grid-cols-1 gap-4 hover:cursor-default"
  >
    <p class="flex justify-center text-xl text-gray-200">
      Create a New Goal/Task
    </p>
    <div class="grid grid-cols-1 gap-4">
      <label class="text-gray-200">Title:</label>
      <input
        type="text"
        v-model="formData.title"
        class="block h-10 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6 focus:outline-none"
      />
    </div>
    <div>
      <label class="text-gray-200">Description:</label>
      <textarea
        v-model="formData.description"
        class="block h-28 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6 focus:outline-none"
      />
    </div>
    <button
      @click="handleSubmit"
      class="bg-blue-400 mt-4 rounded-lg text-gray-300 hover:bg-blue-500 h-10 py-1.5"
    >
      Add Goal
    </button>
  </form>
</template>
