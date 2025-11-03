<script setup lang="ts">
import { useForm } from "@tanstack/vue-form";
import { watchDebounced } from "@vueuse/core";
import { ref } from "vue";
import { toast } from "vue3-toastify";
import { Menu, MenuButton, MenuItems, MenuItem } from "@headlessui/vue";
import {
	useDeleteGoalCategory,
	useUpdateGoalCategory,
} from "@/features/goals/queries";
import type { GoalCategory } from "@/features/goals/schemas/goal.schema";
import { editGoalCategoryFormSchema } from "@/features/goals/schemas/goal-form.schema";
import { Box, Text, InputField } from "@/shared/components/ui";
import { ModalForm } from "@/shared/components/modals";
import GoalCard from "./GoalCard.vue";
import CreateGoalButton from "./CreateGoalButton.vue";
import CreateGoalForm from "../forms/CreateGoalForm.vue";

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
  <Box flex-direction="col" padding="p-4">
    <header class="flex justify-between">
      <Box flex-direction="col">
        <!-- Title Field -->
        <form.Field name="title">
          <template v-slot="{ field }">
            <InputField
              :id="field.name"
              :name="field.name"
              :value="field.state.value"
              class="text-xl"
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

        <!-- XP Per Goal Field -->
        <form.Field name="xp_per_goal">
          <template v-slot="{ field }">
            <InputField
              as="input"
              :id="field.name"
              :text-color="
                field.state.meta.errors.length > 0 ? 'error' : 'light'
              "
              :name="field.name"
              :value="field.state.value"
              type="number"
              class="items-center text-xs"
              container-width="w-full"
              width="w-1/2"
              compact
              @input="
                (value: number | string | null) => {
                  if (typeof value !== 'number') return;
                  field.handleChange(value);
                  console.log(field.state.meta.errors.length);
                }
              "
              @blur="field.handleBlur"
            >
              <template #left>
                <Text as="span" size="xs" color="light">Earn</Text>
              </template>
              <template #right>
                <span class="text-xs text-gray-300">xp/goal</span>
              </template>
            </InputField>
          </template>
        </form.Field>
      </Box>
      <div class="flex">
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
        <Menu as="div" class="relative inline-block">
          <MenuButton class="hover:cursor-pointer">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              class="size-6 stroke-gray-300"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M8.625 12a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H8.25m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H12m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0h-.375M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
              />
            </svg>
          </MenuButton>
          <MenuItems
            class="absolute flex flex-col items-start w-56 bg-gray-500 p-1 rounded-md justify-self-start right-0"
          >
            <MenuItem
              as="button"
              class="w-full px-2 flex justify-start gap-x-2 bg-gray-600 hover:cursor-pointer hover:bg-gray-400 text-gray-700"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                class="size-5 stroke-gray-300"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L6.832 19.82a4.5 4.5 0 0 1-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 0 1 1.13-1.897L16.863 4.487Zm0 0L19.5 7.125"
                />
              </svg>
              <span class="text-sm text-gray-300">Edit</span>
            </MenuItem>

            <MenuItem
              as="button"
              class="w-full px-2 hover:cursor-pointer hover:bg-gray-400 bg-gray-600 text-gray-700 flex gap-x-2"
              @click="handleDeleteCategory"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                class="size-5 stroke-gray-300"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
                />
              </svg>
              <span class="text-sm text-gray-300">Delete</span>
            </MenuItem>
          </MenuItems>
        </Menu>
      </div>
    </header>
    <Box flex-direction="col" gap="gap-4" v-for="goal in goalCategory.goals">
      <GoalCard :goal="goal" />
    </Box>
  </Box>
</template>
<style scoped>
/* Hide default increment and decrement arrows */
.num-input::-webkit-outer-spin-button,
.num-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.num-input {
  -moz-appearance: textfield; /* Hides the arrows in Firefox */
  appearance: textfield; /* Hides the arrows in other browsers */
}
</style>
