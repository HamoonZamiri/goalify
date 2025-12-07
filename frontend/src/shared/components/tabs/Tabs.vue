<script setup lang="ts">
import { TabGroup, TabList, TabPanels } from "@headlessui/vue";
import { computed } from "vue";

const props = withDefaults(
	defineProps<{
		modelValue?: number;
		variant?: "bordered" | "lifted" | "boxed";
	}>(),
	{
		variant: "boxed",
	},
);

const emit = defineEmits<{
	"update:modelValue": [value: number];
}>();

const tabListClasses = computed(() => [
	"tabs",
	"flex",
	"flex-row",
	props.variant === "bordered" && "tabs-bordered",
	props.variant === "lifted" && "tabs-lifted",
	props.variant === "boxed" && "tabs-boxed",
]);
</script>

<template>
	<TabGroup
		:selected-index="props.modelValue"
		@change="(index: number) => emit('update:modelValue', index)"
	>
		<TabList role="tablist" :class="tabListClasses">
			<slot name="tabs"/>
		</TabList>

		<TabPanels v-if="$slots.panels">
			<slot name="panels"/>
		</TabPanels>
	</TabGroup>
</template>
