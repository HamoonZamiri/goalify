<script setup lang="ts">
import { ref } from "vue";
import { useCreateGoal } from "@/features/goals/queries";
import type { CreateGoalFormData } from "@/features/goals/schemas/goal-form.schema";
import { Button, InputField } from "@/shared/components/ui";

const formData = ref<CreateGoalFormData>({
  title: "",
  description: "",
  category_id: "",
});

const props = defineProps<{
  categoryId: string;
}>();

const emit = defineEmits(["submit", "close"]);

const { mutate: createGoal, isPending } = useCreateGoal();

function submit() {
  const submitData: CreateGoalFormData = {
    ...formData.value,
    category_id: props.categoryId,
  };
  emit("submit", { ...submitData });

  createGoal(submitData, {
    onSuccess: () => {
      formData.value.title = "";
      formData.value.description = "";
      emit("close");
    },
  });
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
    </InputField>
    <InputField
      bg="primary"
      text-color="dark"
      type="text"
      v-model="formData.description"
      as="textarea"
    >
      <template #label><Text>Description</Text></template>
    </InputField>
    <Button type="submit" :disabled="isPending">
      <Text>{{ isPending ? "Creating..." : "Add Goal" }}</Text>
    </Button>
  </form>
</template>
