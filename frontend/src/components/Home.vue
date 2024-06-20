<script setup lang="ts">
import GoalCategory from "./goals/GoalCategory.vue";
import { ref } from "vue";
import { type TGoalCategory, mockGoalCategory } from "@/types";
const categories = ref<TGoalCategory[]>([mockGoalCategory, mockGoalCategory]);
</script>

<template>
  <section class="flex gap-4 px-6">
    <div class="p-4 w-1/4" v-for="cat in categories">
      <GoalCategory :goalCategory="cat" />
    </div>
  </section>
</template>
