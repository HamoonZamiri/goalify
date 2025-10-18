<script setup lang="ts">
import { computed } from "vue";
import {
	type ButtonVariant,
	buttonVariantClasses,
	type Height,
	type Width,
} from "@/utils/tailwind";

const props = withDefaults(
	defineProps<{
		variant?: ButtonVariant;
		height?: Height;
		width?: Width;
		class?: string;
	}>(),
	{
		variant: "primary",
		height: "h-10",
		width: "w-full",
	},
);

const classes = computed(() => [
	"inline-flex justify-center items-center gap-2 rounded-lg text-sm font-medium transition-colors",
	buttonVariantClasses[props.variant] ?? "",
	props.height,
	props.width,
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
