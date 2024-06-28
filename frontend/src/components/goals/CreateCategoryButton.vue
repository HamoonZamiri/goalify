<script setup lang="ts">
const props = defineProps<{
  setIsOpen: (value: boolean) => void;
  isOpen: boolean;
}>();
</script>
<template>
  <v-icon name="co-plus" />
  <span
    class="text-xl text-blue-400 font-semibold"
    @click="props.setIsOpen(true)"
  >
    Add Category
  </span>
</template>
