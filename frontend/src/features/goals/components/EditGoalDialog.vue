<script setup lang="ts">
import { ref } from "vue";
import EditGoalForm from "@/features/goals/forms/EditGoalForm.vue";
import type { Goal } from "@/features/goals/schemas/goal.schema";
import { Dialog } from "@/shared/components/ui";

const props = defineProps<{
	goal: Goal;
	modelValue: boolean;
}>();

const emit = defineEmits<{
	"update:modelValue": [value: boolean];
}>();

const editFormRef = ref<InstanceType<typeof EditGoalForm>>();

async function handleClose() {
	await editFormRef.value?.saveIfDirty();
	emit("update:modelValue", false);
}
</script>

<template>
	<Dialog
		:modelValue="modelValue"
		@update:modelValue="handleClose"
		size="drawer"
		position="sideDrawer"
		title="Edit Goal"
		description="Update goal details and status"
	>
		<EditGoalForm ref="editFormRef" :goal="props.goal" @close="handleClose"/>
	</Dialog>
</template>
