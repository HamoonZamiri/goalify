<script setup lang="ts">
import { ref, watch } from "vue";
import { TransitionRoot, Dialog, DialogPanel } from "@headlessui/vue";
import { toast } from "vue3-toastify";
import EditGoalForm from "@/features/goals/forms/EditGoalForm.vue";
import type { Goal } from "@/features/goals/schemas/goal.schema";
import { useUpdateGoal } from "@/features/goals/queries";
import { Text } from "@/shared/components/ui";
import { CheckOutline } from "@/shared/components/icons";

const props = defineProps<{
	goal: Goal;
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
  <header
    @click="openEditingDialog()"
    class="flex p-4 w-full h-full bg-gray-700 hover:cursor-pointer hover:bg-gray-600 gap-x-2 items-center rounded-sm"
    :class="{
      'bg-green-600 hover:bg-green-700': currentStatus === 'complete',
    }"
  >
    <CheckOutline :onclick="handleCheckClick" class="hover:cursor-pointer" />
    <Text as="span" weight="semibold">{{ props.goal.title }}</Text>
  </header>
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
