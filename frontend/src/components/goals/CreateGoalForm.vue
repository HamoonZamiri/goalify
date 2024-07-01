<script setup lang="ts">
import goalState from "@/state/goals";
import { ApiClient } from "@/utils/api";
import { ref } from "vue";

type CreateGoalForm = {
  title: string;
  description: string;
};

const formData = ref<CreateGoalForm>({
  title: "",
  description: "",
});

const error = ref<string | null>(null);

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
  const res = await ApiClient.createGoal(
    title,
    description,
    CreateGoalFormProps.props.categoryId,
  );
  if (typeof res === "string") {
    error.value = res;
    return;
  }
  goalState.addGoal(CreateGoalFormProps.props.categoryId, res.data);
  formData.value.title = "";
  formData.value.description = "";
  CreateGoalFormProps.setIsOpen(false);
}
</script>

<template>
  <form
    class="rounded-lg border bg-white p-3 w-[400px] grid grid-cols-1 gap-8 hover:cursor-default"
  >
    <p class="font-semibold">Create a New Goal/Task</p>
    <div class="grid grid-cols-1 gap-4">
      <label>Title:</label>
      <input
        type="text"
        v-model="formData.title"
        class="block w-full rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-200 sm:text-sm sm:leading-6"
      />
    </div>
    <div>
      <label>Description:</label>
      <textarea
        v-model="formData.description"
        class="block w-full rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-200 sm:text-sm sm:leading-6"
      />
    </div>
    <button
      @click="handleSubmit"
      class="block w-full hover:bg-blue-200 bg-blue-100 rounded-lg h-10"
    >
      Add Goal
    </button>
  </form>
</template>
