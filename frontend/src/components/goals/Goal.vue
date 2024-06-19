<script setup lang="ts">
import { type TGoal } from "@/types";

const props = defineProps<{
  goal: TGoal;
}>();
</script>
<template>
  <div>
    <header class="flex gap-2">
      <span>Check</span>
      <span>{{ props.goal.title }}</span>
    </header>
    <section>
      <p>{{ props.goal.created_at.toDateString() }}</p>
    </section>
  </div>
</template>
