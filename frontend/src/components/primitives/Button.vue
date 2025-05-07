<script setup lang="ts">
import { buttonVariantClasses, type ButtonVariant } from "@/utils/tailwind";
import { computed } from "vue";

const props = defineProps<{
  variant: ButtonVariant;
  label?: string;
  class?: string;
}>();

withDefaults(props, {
  variant: "primary",
});

const classes = computed(() => [
  "inline-flex items-center gap-2 rounded-md text-sm font-medium transition-colors",
  buttonVariantClasses[props.variant] ?? "",
  props.class ?? "",
]);
</script>

<template>
  <button :class="classes">
    <slot name="left" />
    <slot />
    <slot name="right" />
  </button>
</template>
