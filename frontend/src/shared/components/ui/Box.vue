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

const bgClasses = {
	darkest: "bg-gray-900",
	dark: "bg-gray-800",
};

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
		bg?: keyof typeof bgClasses;
	}>(),
	{
		rounded: "rounded-xl",
		shadow: "shadow-sm",
		bordered: false,
		flexDirection: "col",
		bg: "dark",
		class: "",
	},
);
</script>

<template>
  <div
    :class="[
      'text-gray-300 flex',
      flexDirectionClasses[flexDirection],
      bgClasses[bg],
      gap,
      padding,
      rounded,
      shadow,
      bordered ? 'border border-gray-700' : '',
      height,
      width,
      props.class,
    ]"
    v-bind="$attrs"
  >
    <slot />
  </div>
</template>
