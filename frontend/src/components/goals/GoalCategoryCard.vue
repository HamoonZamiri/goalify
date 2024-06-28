<script setup lang="ts">
import type { GoalCategory } from "@/utils/schemas";
import GoalCard from "./GoalCard.vue";
const props = defineProps<{
  goalCategory: GoalCategory;
}>();
</script>
<template>
  <div class="flex-col">
    <header class="flex justify-between">
      <span class="text-xl text-blue-400 font-semibold pb-2">{{
        props.goalCategory.title
      }}</span>
      <v-icon class="hover:cursor-pointer" name="px-notes-plus" />
    </header>
    <section class="">
      <div
        class="rounded-xl border-4 border-blue-400 p-4 mb-4"
        v-for="goal in goalCategory.goals"
      >
        <GoalCard :goal="goal" />
      </div>
    </section>
  </div>
</template>
