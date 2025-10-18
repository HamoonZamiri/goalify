<script setup lang="ts">
import useApi from "@/hooks/api/useApi";
import useGoals from "@/hooks/goals/useGoals";
import type { ErrorResponse } from "@/utils/schemas";
import Text from "@/components/primitives/Text.vue";
import InputField from "@/components/primitives/InputField.vue";
import Button from "@/components/primitives/Button.vue";
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
    class="rounded-lg border bg-gray-800 p-10 w-[95vw] sm:w-[40vw] flex flex-col gap-4 hover:cursor-default"
  >
    <Text as="p" size="xl" class="text-center">
      Create a New Goal/Task Category
    </Text>
    <InputField
      bg="primary"
      text-color="dark"
      type="text"
      v-model="formData.title"
    >
      <template #label><Text>Title</Text></template>
      <template v-if="error?.errors?.title" #error>
        <Text color="error">{{ error?.errors?.title }}</Text>
      </template>
    </InputField>
    <InputField
      bg="primary"
      text-color="dark"
      type="number"
      v-model="formData.xp_per_goal"
    >
      <template #label><Text>XP/Goal:</Text></template>
      <template v-if="error?.errors?.xp_per_goal" #error>
        <Text color="error">{{ error?.errors?.xp_per_goal }}</Text>
      </template>
    </InputField>
    <Button type="submit">
      <Text>Add Category</Text>
    </Button>
  </form>
</template>
