<script setup lang="ts">
import { ref, watch } from "vue";
import { TransitionRoot, Dialog, DialogPanel } from "@headlessui/vue";
import { toast } from "vue3-toastify";
import EditGoalForm from "@/features/goals/forms/EditGoalForm.vue";
import type { Goal } from "@/features/goals/schemas/goal.schema";
import { useUpdateGoal } from "@/features/goals/queries";
import { Text, Box } from "@/shared/components/ui";
import { CheckOutline } from "@/shared/components/icons";

const props = defineProps<{
	goal: Goal;
	xpPerGoal: number;
}>();

const { mutateAsync: updateGoal } = useUpdateGoal();

const isEditing = ref(false);
const editFormRef = ref<InstanceType<typeof EditGoalForm>>();

function setIsEditing(value: boolean) {
	isEditing.value = value;
}

function openEditingDialog() {
	setIsEditing(true);
}

async function handleClose() {
	await editFormRef.value?.saveIfDirty();
	setIsEditing(false);
}

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
	<Box
		data-testid="goal-card"
		flex-direction="row"
		gap="gap-4"
		@click.stop="() => openEditingDialog()"
		class="hover:cursor-pointer hover:bg-gray-700 items-center justify-between p-1 rounded-xl"
	>
		<Box flex-direction="row" class="gap-x-2 items-center" bg="inherit">
			<CheckOutline
				:onclick="handleCheckClick"
				class="hover:cursor-pointer"
				:fill="currentStatus === 'complete' ? 'green' : 'none'"
			/>
			<Text
				as="span"
				size="sm"
				weight="normal"
				:class="currentStatus === 'complete' ? 'line-through opacity-50' : ''"
			>
				{{ props.goal.title }}
			</Text>
		</Box>
		<Text as="span" size="sm" weight="normal">{{`${props.xpPerGoal} XP`}}</Text>
	</Box>
	<section>
		<TransitionRoot
			:show="isEditing"
			appear
			enter="transition-all ease-in-out duration-500 transform"
		>
			<Dialog
				class="absolute inset-0 h-screen flex justify-end hover:cursor-pointer z-10 w-screen bg-opacity-10 rounded-lg"
				@close="handleClose"
			>
				<DialogPanel class="w-full sm:w-1/2">
					<EditGoalForm
						ref="editFormRef"
						:goal="props.goal"
						@close="handleClose"
					/>
				</DialogPanel>
			</Dialog>
		</TransitionRoot>
	</section>
</template>
