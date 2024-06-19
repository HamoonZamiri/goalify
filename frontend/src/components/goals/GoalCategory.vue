<script setup lang="ts">
import Goal from "./Goal.vue";
import { type TGoalCategory } from "@/types";
const props = defineProps<{
  goalCategory: TGoalCategory;
}>();
</script>
<template>
  <div class="flex-col">
    <header class="flex justify-between">
      <span class="text-xl text-blue-400 font-semibold pb-2">{{
        props.goalCategory.title
      }}</span>
      <span>+</span>
    </header>
    <section class="">
      <div
        class="bg-blue-100 p-2 my-4"
        v-for="goal in props.goalCategory.goals"
      >
        <Goal :goal="goal" />
      </div>
    </section>
  </div>
</template>
