<script setup lang="ts">
import { useForm } from "@tanstack/vue-form";
import { toast } from "vue3-toastify";
import { useCreateGoal } from "@/features/goals/queries";
import {
	type CreateGoalFormData,
	createGoalFormSchema,
} from "@/features/goals/schemas/goal-form.schema";
import { Text, InputField, Button } from "@/shared/components/ui";
import { ArrowPath } from "@/shared/components/icons";

const props = defineProps<{
	categoryId: string;
}>();

const emit = defineEmits(["submit", "close"]);

const { mutateAsync: createGoal, isPending } = useCreateGoal();

async function handleSubmit(data: CreateGoalFormData) {
	try {
		const result = await createGoal(data);
		toast.success(`Successfully created goal: ${result.title}`);
		emit("submit", result);
		emit("close");
	} catch (error) {
		toast.error(
			`Failed to create goal: ${error instanceof Error ? error.message : "Unknown error"}`,
		);
	}
}

const form = useForm({
	defaultValues: {
		title: "",
		description: "",
		category_id: props.categoryId,
	},
	validators: {
		onChange: createGoalFormSchema,
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
		<Text as="p" size="xl" class="text-center">Create a New Goal/Task </Text>
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
		<form.Field name="description">
			<template v-slot="{ field }">
				<InputField
					:id="field.name"
					:name="field.name"
					:value="field.state.value"
					bg="primary"
					text-color="dark"
					type="text"
					as="textarea"
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
						<Text>Description</Text>
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
					<Text v-else>Add Goal </Text>
				</Button>
			</template>
		</form.Subscribe>
	</form>
</template>
