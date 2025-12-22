<script setup lang="ts">
import CreateGoalForm from "@/features/goals/forms/CreateGoalForm.vue";
import { Dialog } from "@/shared/components/ui";

defineProps<{
	categoryId: string;
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
		title="Create Goal"
		description="Add a new goal to this category"
	>
		<CreateGoalForm :category-id="categoryId" @close="handleClose"/>
	</Dialog>
</template>
