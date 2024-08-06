<script setup lang="ts">
const props = defineProps<{
  setIsOpen: (value: boolean) => void;
  isOpen: boolean;
}>();
</script>
<template>
  <v-icon name="co-plus" fill="white" />
  <span class="text-xl text-gray-200" @click="props.setIsOpen(true)">
    Add Category
  </span>
</template>
