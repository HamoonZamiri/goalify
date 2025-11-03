<script setup lang="ts">
import { useForm } from "@tanstack/vue-form";
import { watchDebounced } from "@vueuse/core";
import { ref } from "vue";
import { toast } from "vue3-toastify";
import {
  Listbox,
  ListboxButton,
  ListboxOptions,
  ListboxOption,
} from "@headlessui/vue";
import { useDeleteGoal, useUpdateGoal } from "@/features/goals/queries";
import type { Goal } from "@/features/goals/schemas/goal.schema";
import { editGoalFormSchema } from "@/features/goals/schemas/goal-form.schema";
import { Box, Text, InputField, Button } from "@/shared/components/ui";
import { XMark } from "@/shared/components/icons";
import { DeleteModal } from "@/shared/components/modals";

const props = defineProps<{
  goal: Goal;
}>();

const emit = defineEmits<{
  close: [];
  update: [goal: Goal];
}>();

const isDeleting = ref(false);

const { mutateAsync: updateGoal } = useUpdateGoal();
const { mutateAsync: deleteGoalMutation } = useDeleteGoal();

const form = useForm({
  defaultValues: {
    title: props.goal.title,
    description: props.goal.description,
    status: props.goal.status,
  },
  validators: {
    onChange: editGoalFormSchema,
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
      await updateGoal({
        goalId: props.goal.id,
        data: {
          title: values.title,
          description: values.description,
          status: values.status,
        },
      });
    } catch (error) {
      toast.error(
        `Failed to update goal: ${error instanceof Error ? error.message : "Unknown error"}`,
      );
    }
  },
  { debounce: 500, deep: true },
);

async function saveIfDirty() {
  if (isDirty.value && isValid.value) {
    try {
      await updateGoal({
        goalId: props.goal.id,
        data: {
          title: formValues.value.title,
          description: formValues.value.description,
          status: formValues.value.status,
        },
      });
    } catch (error) {
      toast.error(
        `Failed to save goal: ${error instanceof Error ? error.message : "Unknown error"}`,
      );
    }
  }
}

async function handleDeleteGoal(e: MouseEvent) {
  e.preventDefault();
  try {
    await deleteGoalMutation(props.goal.id);
    toast.success("Successfully deleted goal");
    emit("close");
  } catch (error) {
    toast.error(
      `Failed to delete goal: ${error instanceof Error ? error.message : "Unknown error"}`,
    );
  }
}

function handleClose() {
  emit("close");
}

const statuses = [
  { id: 1, name: "Not Complete", value: "not_complete" },
  { id: 2, name: "Complete", value: "complete" },
];

const statusMap = {
  not_complete: "Not Complete",
  complete: "Complete",
} as const;

defineExpose({ saveIfDirty });
</script>

<template>
  <Box
    gap="gap-4"
    shadow="shadow-lg"
    flex-direction="col"
    padding="p-8"
    height="h-full"
    width="w-full"
    class="border-white hover:cursor-default shadow-gray-400"
  >
    <!-- Title Field -->
    <form.Field name="title">
      <template v-slot="{ field }">
        <InputField
          :id="field.name"
          :name="field.name"
          :value="field.state.value"
          type="text"
          containerWidth="w-full"
          class="text-3xl text-gray-300"
          @input="
            (value: string | number | null) => {
              if (typeof value !== 'string') return;
              field.handleChange(value);
            }
          "
          @blur="field.handleBlur"
        >
          <template #right>
            <XMark
              :onclick="handleClose"
              class="sm:ml-auto hover:cursor-pointer"
            />
          </template>
        </InputField>
      </template>
    </form.Field>

    <!-- Status Field -->
    <Box flex-direction="col">
      <Text size="xl" as="p">Status</Text>
      <form.Field name="status">
        <template v-slot="{ field }">
          <Box flex-direction="col" gap="gap-2">
            <Listbox
              as="div"
              class="relative"
              :value="field.state.value"
              @update:model-value="field.handleChange"
              :model-value="field.state.value"
            >
              <ListboxButton
                :class="{
                  'w-56 h-8 text-center rounded-lg text-gray-600': true,
                  'bg-green-400': field.state.value === 'complete',
                  'bg-orange-400': field.state.value !== 'complete',
                }"
              >
                <Text color="dark">{{ statusMap[field.state.value] }}</Text>
              </ListboxButton>
              <ListboxOptions class="absolute z-10 mt-1 flex flex-col gap-1">
                <ListboxOption
                  v-for="status in statuses"
                  :key="status.id"
                  :value="status.value"
                  :disabled="status.value === field.state.value"
                  :class="{
                    'w-56 hover:cursor-pointer h-8 bg-gray-300 text-gray-600 text-center hover:bg-gray-400 rounded-lg inline-flex items-center justify-center': true,
                    hidden: field.state.value === status.value,
                  }"
                >
                  <Text color="dark">{{ status.name }}</Text>
                </ListboxOption>
              </ListboxOptions>
            </Listbox>
          </Box>
        </template>
      </form.Field>
    </Box>

    <!-- Description Field -->
    <Box flex-direction="col">
      <Text size="xl" as="p">Description</Text>
      <form.Field name="description">
        <template v-slot="{ field }">
          <textarea
            :id="field.name"
            :name="field.name"
            :value="field.state.value"
            @input="
              (e) => {
                field.handleChange((e.target as HTMLTextAreaElement).value);
              }
            "
            @blur="field.handleBlur"
            class="w-full bg-gray-300 focus:outline-none h-64 p-2 text-gray-600"
          />
        </template>
      </form.Field>
    </Box>

    <!-- Delete Button -->
    <Button
      variant="secondary"
      class="hover:bg-red-600 h-10"
      @click="isDeleting = true"
    >
      <Text weight="semibold" size="base">Delete This Goal</Text>
    </Button>

    <DeleteModal
      :is-open="isDeleting"
      :set-opener="(val: boolean) => (isDeleting = val)"
      delete-message="Are you sure you want to delete this goal?"
      :delete-function="handleDeleteGoal"
    />
  </Box>
</template>
