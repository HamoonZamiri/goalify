<script setup lang="ts">
const props = defineProps<{
  setIsOpen: (value: boolean) => void;
  isOpen: boolean;
}>();
</script>
<template>
  <v-icon name="co-plus" fill="white" />
  <span class="text-xl text-white font-semibold" @click="props.setIsOpen(true)">
    Add Category
  </span>
</template>
