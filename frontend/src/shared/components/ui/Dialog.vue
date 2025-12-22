<script setup lang="ts">
import {
	DialogDescription,
	DialogPanel,
	DialogTitle,
	Dialog as HeadlessDialog,
	TransitionRoot,
} from "@headlessui/vue";
import {
	type DialogPattern,
	type DialogSize,
	dialogPatterns,
} from "@/utils/tailwind";
import Text from "./Text.vue";

const props = withDefaults(
	defineProps<{
		modelValue: boolean;
		size?: DialogSize;
		position?: DialogPattern;
		title?: string;
		description?: string;
	}>(),
	{
		size: "md",
		position: "centered",
	},
);

const emit = defineEmits<{
	"update:modelValue": [value: boolean];
}>();

function closeDialog() {
	emit("update:modelValue", false);
}
</script>

<template>
	<TransitionRoot
		:show="modelValue"
		appear
		enter="transition-opacity duration-200"
		enterFrom="opacity-0"
		enterTo="opacity-100"
		leave="transition-opacity duration-150"
		leaveFrom="opacity-100"
		leaveTo="opacity-0"
	>
		<HeadlessDialog
			:open="modelValue"
			@close="closeDialog"
			class="relative z-50"
		>
			<!-- Backdrop -->
			<div :class="dialogPatterns.backdrop" aria-hidden="true"/>

			<!-- Container with positioning -->
			<div :class="`${dialogPatterns.base} ${dialogPatterns[position]}`">
				<!-- Panel -->
				<DialogPanel :class="dialogPatterns.sizes[size]">
					<DialogTitle v-if="title" as="h2" class="sr-only">
						<Text as="span">{{ title }}</Text>
					</DialogTitle>

					<!-- Optional description (for accessibility) -->
					<DialogDescription v-if="description" as="p" class="sr-only">
						<Text as="span">{{ description }}</Text>
					</DialogDescription>

					<!-- Content slot -->
					<slot/>
				</DialogPanel>
			</div>
		</HeadlessDialog>
	</TransitionRoot>
</template>
