<script setup lang="ts">
import { useForm } from "@tanstack/vue-form";
import { toast } from "vue3-toastify";
import { useCreateGoalCategory } from "@/features/goals/queries";
import {
	type CreateGoalCategoryFormData,
	createGoalCategoryFormSchema,
} from "@/features/goals/schemas/goal-form.schema";
import { Text, InputField, Button } from "@/shared/components/ui";
import { ArrowPath } from "@/shared/components/icons";

const emit = defineEmits(["submit", "close"]);

const { mutateAsync: createCategory, isPending } = useCreateGoalCategory();

async function handleSubmit(data: CreateGoalCategoryFormData) {
	try {
		const result = await createCategory(data);
		emit("submit", result);
		emit("close");
	} catch (error) {
		toast.error(
			`Failed to create category: ${error instanceof Error ? error.message : "Unknown error"}`,
		);
	}
}

const form = useForm({
	defaultValues: {
		title: "",
		xp_per_goal: 10,
	},
	validators: {
		onChange: createGoalCategoryFormSchema,
	},
	onSubmit: async ({ value }) => {
		await handleSubmit(value);
	},
});
</script>

<template>
	<form
		@submit="
      (e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }
    "
		class="rounded-lg border bg-gray-800 p-6 w-[95vw] sm:w-[40vw] flex flex-col space-y-4 hover:cursor-default"
	>
		<Text as="p" size="xl" class="text-center">Create a New Category </Text>
		<form.Field name="title">
			<template v-slot="{ field, state }">
				<InputField
					:id="field.name"
					:name="field.name"
					:value="field.state.value"
					bg="primary"
					text-color="dark"
					type="text"
					@input="
            (value: string | number | null) => {
              if (typeof value !== 'string') return;
              field.handleChange(value);
            }
          "
					@blur="field.handleBlur"
					errorslot
				>
					<template #label>
						<Text>Title</Text>
					</template>
					<template
						#error
						v-if="
              field.state.meta.isTouched && field.state.meta.errors.length > 0
            "
					>
						<Text color="error">{{ field.state.meta.errors[0]?.message }}</Text>
					</template>
				</InputField>
			</template>
		</form.Field>
		<form.Field name="xp_per_goal">
			<template v-slot="{ field, state }">
				<InputField
					:id="field.name"
					:name="field.name"
					:value="field.state.value"
					bg="primary"
					text-color="dark"
					type="number"
					@input="
            (value: string | number | null) => {
              if (typeof value !== 'number') return;
              field.handleChange(value);
            }
          "
					@blur="field.handleBlur"
					errorslot
				>
					<template #label>
						<Text>XP Per Goal</Text>
					</template>
					<template
						#error
						v-if="
              field.state.meta.isTouched && field.state.meta.errors.length > 0
            "
					>
						<Text color="error">{{ field.state.meta.errors[0]?.message }}</Text>
					</template>
				</InputField>
			</template>
		</form.Field>
		<form.Subscribe>
			<template v-slot="{ canSubmit, isSubmitting }">
				<Button
					type="submit"
					:disabled="!canSubmit || isPending || isSubmitting"
				>
					<ArrowPath class="animate-spin" v-if="isSubmitting"/>
					<Text v-else>Add Category </Text>
				</Button>
			</template>
		</form.Subscribe>
	</form>
</template>
