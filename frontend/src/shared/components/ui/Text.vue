<script setup lang="ts">
const props = withDefaults(
	defineProps<{
		as?: keyof HTMLElementTagNameMap;
		size?: "xs" | "sm" | "base" | "lg" | "xl" | "2xl" | "3xl";
		weight?: "normal" | "medium" | "bold" | "semibold";
		color?: "light" | "dark" | "error";
		class?: string;
	}>(),
	{
		as: "span",
		size: "base",
		weight: "normal",
		color: "light",
	},
);

const colorClasses = {
	light: "text-gray-300",
	dark: "text-gray-600",
	error: "text-red-400",
};

const sizeClasses = {
	xs: "text-xs",
	sm: "text-sm",
	base: "text-base",
	lg: "text-lg",
	xl: "text-xl",
	"2xl": "text-2xl",
	"3xl": "text-3xl",
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
      colorClasses[color],
    ]"
  >
    <slot />
  </component>
</template>
