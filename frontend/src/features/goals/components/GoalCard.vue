<script setup lang="ts">
import { ref, watch } from "vue";
import { toast } from "vue3-toastify";
import EditGoalDialog from "@/features/goals/components/EditGoalDialog.vue";
import type { Goal } from "@/features/goals/schemas/goal.schema";
import { useUpdateGoal } from "@/features/goals/queries";
import { Text, ClickSurface, OverlayButton } from "@/shared/components/ui";
import { IconButton } from "@/shared/components/icons";

const props = defineProps<{
	goal: Goal;
}>();

const { mutateAsync: updateGoal } = useUpdateGoal();

const isEditing = ref(false);

async function handleCheckClick() {
	const newStatus =
		currentStatus.value === "complete" ? "not_complete" : "complete";

	try {
		await updateGoal({
			goalId: props.goal.id,
			data: { status: newStatus },
		});
	} catch (error) {
		toast.error(
			`Failed to update goal status: ${error instanceof Error ? error.message : "Unknown error"}`,
		);
	}
}

/**
 * Watch for goal status changes to update the card appearance
 */
const currentStatus = ref(props.goal.status);
watch(
	() => props.goal,
	(newGoal) => {
		currentStatus.value = newGoal.status;
	},
	{ deep: true },
);
</script>
<template>
	<ClickSurface
		data-testid="goal-card"
		class="flex flex-row items-center justify-between gap-4 p-1 rounded-xl hover:bg-gray-700"
	>
		<template #overlay>
			<OverlayButton aria-label="Edit goal" @click.stop="isEditing = true" />
		</template>
		<Text
			as="span"
			size="sm"
			weight="normal"
			:class="currentStatus === 'complete' ? 'line-through opacity-50' : ''"
		>
			{{ props.goal.title }}
		</Text>
		<template #actions>
			<IconButton
				icon="check-outline"
				ariaLabel="Toggle goal completion"
				:icon-fill="currentStatus === 'complete' ? 'green' : 'none'"
				@click.stop="handleCheckClick"
			/>
		</template>
	</ClickSurface>
	<EditGoalDialog v-model="isEditing" :goal="props.goal" />
</template>
