<script setup lang="ts">
import Button from "@/shared/components/ui/Button.vue";
import type {
	ButtonVariant,
	IconFill,
	IconName,
	IconSize,
} from "@/utils/tailwind";
import Icon from "./Icon.vue";

const props = withDefaults(
	defineProps<{
		icon: IconName;
		variant?: ButtonVariant;
		size?: IconSize;
		class?: string;
		disabled?: boolean;
		ariaLabel: string;
		iconFill?: IconFill;
	}>(),
	{
		variant: "ghost",
		size: "md",
		disabled: false,
	},
);

defineEmits<{
	click: [event: MouseEvent];
}>();
</script>

<template>
	<Button
		:variant="variant"
		width="w-auto"
		class="p-0"
		:disabled="disabled"
		:aria-label="ariaLabel"
		@click="$emit('click', $event)"
	>
		<Icon :name="icon" :size="size" :class="props.class" :fill="iconFill"/>
	</Button>
</template>
