<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    as?: keyof HTMLElementTagNameMap;
    size?: "sm" | "base" | "lg";
    weight?: "normal" | "medium" | "bold" | "semibold";
    class?: string;
  }>(),
  {
    as: "span",
    size: "base",
    weight: "normal",
  },
);

const sizeClasses = {
  sm: "text-sm",
  base: "text-base",
  lg: "text-lg",
  xl: "text-xl",
  "2xl": "text-2xl",
};

const weightClasses = {
  normal: "font-normal",
  medium: "font-medium",
  bold: "font-bold",
  semibold: "font-semibold",
};
</script>

<template>
  <component
    :is="as"
    :class="[
      sizeClasses[size],
      weightClasses[weight],
      props.class ?? '',
      'text-gray-300',
    ]"
  >
    <slot />
  </component>
</template>
