<script setup lang="ts">
import { useForm } from "@tanstack/vue-form";
import { watchDebounced } from "@vueuse/core";
import { ref } from "vue";
import { toast } from "vue3-toastify";
import { Disclosure, DisclosureButton, DisclosurePanel } from "@headlessui/vue";
import {
	useDeleteGoalCategory,
	useUpdateGoalCategory,
} from "@/features/goals/queries";
import type { GoalCategory } from "@/features/goals/schemas/goal.schema";
import { editGoalCategoryFormSchema } from "@/features/goals/schemas/goal-form.schema";
import { Box, Text, InputField } from "@/shared/components/ui";
import { ChevronUp } from "@/shared/components/icons";
import { ModalForm } from "@/shared/components/modals";
import { GoalCard, CreateGoalButton, CreateGoalForm } from "@/features/goals";

const props = defineProps<{
	goalCategory: GoalCategory;
}>();

const isCreateGoalDialogOpen = ref(false);

const { mutateAsync: updateCategory } = useUpdateGoalCategory();
const { mutateAsync: deleteCategory } = useDeleteGoalCategory();

const form = useForm({
	defaultValues: {
		title: props.goalCategory.title,
		xp_per_goal: props.goalCategory.xp_per_goal,
	},
	validators: {
		onChange: editGoalCategoryFormSchema,
	},
});

const formValues = form.useStore((state) => state.values);
const isDirty = form.useStore((state) => state.isDirty);
const isValid = form.useStore((state) => state.isValid);

watchDebounced(
	formValues,
	async (values) => {
		if (!isDirty.value || !isValid.value) return;

		try {
			await updateCategory({
				categoryId: props.goalCategory.id,
				data: {
					title: values.title,
					xp_per_goal: values.xp_per_goal,
				},
			});
		} catch (error) {
			toast.error(
				`Failed to update category: ${error instanceof Error ? error.message : "Unknown error"}`,
			);
		}
	},
	{ debounce: 500, deep: true },
);

async function handleDeleteCategory(e: MouseEvent) {
	e.preventDefault();
	try {
		await deleteCategory(props.goalCategory.id);
		toast.success("Successfully deleted category");
	} catch (error) {
		toast.error(
			`Failed to delete category: ${error instanceof Error ? error.message : "Unknown error"}`,
		);
	}
}
</script>

<template>
	<Disclosure as="div" v-slot="{ open }">
		<Box flex-direction="col" padding="p-4">
			<header class="flex justify-between w-full">
				<Box flex-direction="col">
					<!-- Title Field -->
					<form.Field name="title">
						<template v-slot="{ field }">
							<InputField
								:id="field.name"
								:name="field.name"
								:value="field.state.value"
								class="text-md rounded-none"
								@input="
                (value: string | number | null) => {
                  if (typeof value !== 'string') return;
                  field.handleChange(value);
                }
              "
								@blur="field.handleBlur"
							/>
						</template>
					</form.Field>
				</Box>
				<Box class="items-center" flex-direction="row">
					<CreateGoalButton
						class="hover:cursor-pointer"
						@click="isCreateGoalDialogOpen = true"
					/>
					<ModalForm v-model="isCreateGoalDialogOpen">
						<CreateGoalForm
							:category-id="props.goalCategory.id"
							@close="isCreateGoalDialogOpen = false"
						/>
					</ModalForm>
					<DisclosureButton>
						<ChevronUp :class="open ? 'rotate-180 transform' : ''"/>
					</DisclosureButton>
				</Box>
			</header>
			<DisclosurePanel class="transition w-full mt-4">
				<Box flex-direction="col" v-for="goal in goalCategory.goals">
					<GoalCard :goal="goal" :xp-per-goal="goalCategory.xp_per_goal"/>
				</Box>
			</DisclosurePanel>
		</Box>
	</Disclosure>
</template>
