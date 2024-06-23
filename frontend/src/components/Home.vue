<script setup lang="ts">
import GoalCategory from "./goals/GoalCategory.vue";
import { ref } from "vue";
import { type TGoalCategory, mockGoalCategory } from "@/utils/types";
import CreateGoalCategoryDialog from "./goals/CreateGoalCategoryDialog.vue";
const categories = ref<TGoalCategory[]>([mockGoalCategory, mockGoalCategory]);
</script>

<template>
  <div
    class="flex flex-col items-center sm:items-start px-6 w-auto bg-slate-50"
  >
    <section class="flex-col sm:flex-row flex gap-4 w-auto">
      <div class="w-full sm:w-1/3" v-for="cat in categories">
        <GoalCategory :goalCategory="cat" />
      </div>
      <CreateGoalCategoryDialog class="" />
    </section>
  </div>
</template>
