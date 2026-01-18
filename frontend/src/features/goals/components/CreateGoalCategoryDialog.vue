<script setup lang="ts">
import { CreateGoalCategoryForm } from "@/features/goals/forms";
import { Dialog } from "@/shared/components/ui";

defineProps<{
	modelValue: boolean;
}>();

const emit = defineEmits<{
	"update:modelValue": [value: boolean];
}>();

function handleClose() {
	emit("update:modelValue", false);
}
</script>

<template>
	<Dialog
		:modelValue="modelValue"
		@update:modelValue="handleClose"
		size="md"
		position="centered"
		title="Create Goal Category"
		description="Add a new category to organize your goals"
	>
		<CreateGoalCategoryForm @close="handleClose" />
	</Dialog>
</template>
