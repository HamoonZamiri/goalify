<script setup lang="ts">
import useApi from "@/hooks/api/useApi";
import useGoals from "@/hooks/goals/useGoals";
import type { ErrorResponse } from "@/utils/schemas";
import { ref } from "vue";
import { toast } from "vue3-toastify";
import Text from "@/components/primitives/Text.vue";
import InputField from "@/components/primitives/InputField.vue";
import Button from "@/components/primitives/Button.vue";

type CreateGoalForm = {
  title: string;
  description: string;
};

const formData = ref<CreateGoalForm>({
  title: "",
  description: "",
});

const error = ref<ErrorResponse>();
const { addGoal } = useGoals();
const { createGoal, isError } = useApi();

const CreateGoalFormProps = defineProps<{
  categoryId: string;
}>();

const { categoryId } = CreateGoalFormProps;

const emit = defineEmits(["submit", "close"]);

async function submit() {
  emit("submit", { ...formData.value });
  const { title, description } = formData.value;
  const res = await createGoal(title, description, categoryId);
  if (isError(res)) {
    error.value = res;
    return;
  }
  addGoal(categoryId, res);
  formData.value.title = "";
  formData.value.description = "";
  error.value = undefined;

  emit("close");
  toast.success(`Successfully created goal: ${res.title}`);
}
</script>

<template>
  <form
    @submit.prevent="submit"
    class="rounded-lg border bg-gray-800 p-6 w-[95vw] sm:w-[40vw] flex flex-col space-y-4 hover:cursor-default"
  >
    <Text as="p" size="xl" class="text-center"> Create a New Goal/Task </Text>
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
      type="text"
      v-model="formData.description"
      as="textarea"
    >
      <template #label><Text>Description</Text></template>
      <template v-if="error?.errors?.description" #error>
        <Text color="error">{{ error?.errors?.description }}</Text>
      </template>
    </InputField>
    <Button type="submit">
      <Text>Add Goal</Text>
    </Button>
  </form>
</template>
