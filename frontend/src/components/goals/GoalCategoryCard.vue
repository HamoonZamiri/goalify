<script setup lang="ts">
import type { GoalCategory } from "@/utils/schemas";
import GoalCard from "./GoalCard.vue";
import ModalForm from "../ModalForm.vue";
import CreateGoalForm from "./CreateGoalForm.vue";
import CreateGoalButton from "./CreateGoalButton.vue";
const props = defineProps<{
  goalCategory: GoalCategory;
}>();
</script>
<template>
  <div class="flex flex-col">
    <header class="flex justify-between">
      <span class="text-xl text-gray-200 pb-2">{{
        props.goalCategory.title
      }}</span>
      <ModalForm
        :FormComponent="CreateGoalForm"
        :OpenerComponent="CreateGoalButton"
        :formProps="{ categoryId: props.goalCategory.id }"
      />
    </header>
    <div
      class="p-4 mb-4 w-full bg-gray-700 rounded-sm"
      v-for="goal in goalCategory.goals"
    >
      <GoalCard :goal="goal" />
    </div>
  </div>
</template>
