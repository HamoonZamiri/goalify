<script setup lang="ts">
import { ref } from "vue";
import { useCreateGoalCategory } from "@/features/goals/queries";
import type { CreateGoalCategoryFormData } from "@/features/goals/schemas/goal-form.schema";
import { Text, Button, InputField } from "@/shared/components/ui";

const formData = ref<CreateGoalCategoryFormData>({
  title: "",
  xp_per_goal: 10,
});

const emit = defineEmits(["submit", "close"]);

const { mutate: createCategory, isPending } = useCreateGoalCategory();

function submit() {
  emit("submit", { ...formData.value });
  createCategory(formData.value, {
    onSuccess: () => {
      formData.value.title = "";
      formData.value.xp_per_goal = 10;
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
    <Text as="p" size="xl" class="text-center"> Create a New Category </Text>
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
      type="number"
      v-model.number="formData.xp_per_goal"
    >
      <template #label><Text>XP Per Goal</Text></template>
    </InputField>
    <Button type="submit" :disabled="isPending">
      <Text>{{ isPending ? "Creating..." : "Create Category" }}</Text>
    </Button>
  </form>
</template>
