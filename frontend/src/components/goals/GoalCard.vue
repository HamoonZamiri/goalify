<script setup lang="ts">
import type { Goal } from "@/utils/schemas";
const props = defineProps<{
  goal: Goal;
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
  <section></section>
</template>
