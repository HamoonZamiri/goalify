<script setup lang="ts">
import type {
  FlexDirection,
  Gap,
  Height,
  Padding,
  Rounding,
  Shadow,
  Width,
} from "@/utils/tailwind";

const flexDirectionClasses = {
  row: "flex-row",
  col: "flex-col",
} as const;

const props = withDefaults(
  defineProps<{
    padding?: Padding;
    rounded?: Rounding;
    shadow?: Shadow;
    bordered?: boolean;
    flexDirection?: FlexDirection;
    gap?: Gap;
    height?: Height;
    width?: Width;
    class?: string;
  }>(),
  {
    rounded: "rounded-xl",
    shadow: "shadow-sm",
    bordered: false,
    flexDirection: "col",
    class: "",
  },
);
</script>

<template>
  <div
    :class="[
      'bg-gray-800 text-gray-300 flex',
      flexDirectionClasses[flexDirection],
      gap,
      padding,
      rounded,
      shadow,
      bordered ? 'border border-gray-700' : '',
      height,
      props.class,
    ]"
    v-bind="$attrs"
  >
    <slot />
  </div>
</template>
