<script setup lang="ts">
import GoalCategory from "./components/goals/GoalCategory.vue";
import { ref } from "vue";
import { type TGoalCategory, mockGoalCategory } from "./types";
const categories = ref<TGoalCategory[]>([mockGoalCategory, mockGoalCategory]);
</script>

<template>
  <header class="bg-emerald-100 mb-2">
    <div class="flex text-blue-400 justify-between p-6">
      <RouterLink to="/">
        <h1 class="font-semibold text-3xl hover:text-blue-500">Goalify</h1>
      </RouterLink>
      <nav class="flex gap-4">
        <RouterLink
          class="font-semibold text-xl hover:text-blue-500"
          to="/login"
          >Login</RouterLink
        >
        <RouterLink
          class="font-semibold text-xl hover:text-blue-500"
          to="/register"
          >Register</RouterLink
        >
      </nav>
    </div>
  </header>
  <section class="flex gap-4 px-6">
    <div class="p-4 w-1/4" v-for="cat in categories">
      <GoalCategory :goalCategory="cat" />
    </div>
  </section>
  <main><RouterView /></main>
</template>
