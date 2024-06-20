<script setup lang="ts">
import { type TGoal } from "@/types";
const props = defineProps<{
  goal: TGoal;
}>();
</script>
<template>
  <header class="flex gap-x-2 items-center">
    <v-icon
      class="hover:cursor-pointer"
      name="io-checkmark-circle-outline"
      animation="wrench"
      hover
    />
    <span class="font-semibold">{{ props.goal.title }}</span>
  </header>
  <section>
    <p>{{ props.goal.created_at.toDateString() }}</p>
  </section>
</template>
