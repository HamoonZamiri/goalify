<script setup lang="ts">
import { computed } from "vue";
import {
	type ButtonVariant,
	buttonVariantClasses,
	type Height,
	type Width,
	type Padding,
} from "@/utils/tailwind";
import { Icon } from "@/shared/components/icons";

const props = withDefaults(
	defineProps<{
		variant?: ButtonVariant;
		height?: Height;
		width?: Width;
		padding?: Padding;
		class?: string;
		loading?: boolean;
	}>(),
	{
		variant: "primary",
		height: "h-10",
		width: "w-full",
		loading: false,
	},
);

const classes = computed(() => [
	"inline-flex justify-center items-center gap-2 rounded-lg text-sm font-medium transition-colors hover:cursor-pointer",
	buttonVariantClasses[props.variant] ?? "",
	props.height,
	props.width,
	props.padding,
	props.class ?? "",
]);
</script>

<template>
	<button :class="classes">
		<template v-if="loading">
			<Icon name="arrow-path" class="animate-spin" />
		</template>
		<template v-else>
			<slot name="left" />
			<slot />
			<slot name="right" />
		</template>
	</button>
</template>
