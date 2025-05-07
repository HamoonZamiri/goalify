<script setup lang="ts">
import type { Goal } from "@/utils/schemas";
import { ref, watch, reactive, nextTick } from "vue";
import {
  Dialog,
  DialogPanel,
  Listbox,
  ListboxButton,
  ListboxOption,
  ListboxOptions,
  TransitionRoot,
} from "@headlessui/vue";
import Text from "@/components/primitives/Text.vue";
import Box from "@/components/primitives/Box.vue";
import InputField from "@/components/primitives/InputField.vue";
import CheckOutline from "@/components/icons/CheckOutline.vue";
import Button from "@/components/primitives/Button.vue";
import XMark from "@/components/icons/XMark.vue";
("@/components/icons/CheckOutline.vue");
import useGoals from "@/hooks/goals/useGoals";
import useApi from "@/hooks/api/useApi";
import DeleteModal from "@/components/DeleteModal.vue";
import { toast } from "vue3-toastify";
const props = defineProps<{
  goal: Goal;
}>();

const updates = reactive<{
  title: string;
  description: string;
  status: "complete" | "not_complete";
}>({
  title: props.goal.title,
  description: props.goal.description,
  status: props.goal.status,
});

const isEditing = ref(false);
const isDeleting = ref(false);

const { deleteGoal } = useGoals();
const { deleteGoal: deleteGoalApi, updateGoal: updateGoalApi } = useApi();

function setIsEditing(value: boolean) {
  isEditing.value = value;
}

function setIsDeleting(value: boolean) {
  isDeleting.value = value;
}

function openEditingDialog() {
  setIsEditing(false);
  nextTick(() => {
    setIsEditing(true);
  });
}

async function handleDeleteGoal(e: MouseEvent) {
  e.preventDefault();

  await deleteGoalApi(props.goal.id);

  // remove goal from state
  deleteGoal(props.goal.category_id, props.goal.id);

  toast.success(`Successfully deleted goal: ${props.goal.title}`, {
    theme: "dark",
  });

  setIsEditing(false);
  setIsDeleting(false);
}

async function handleToggleStatus(e: MouseEvent) {
  e.preventDefault();
  updates.status = updates.status === "complete" ? "not_complete" : "complete";
}

// watches for updates on the goal title and description
watch(updates, async (newUpdates) => {
  // syncronhize state passed in from props with local reactive updates
  props.goal.title = newUpdates.title;
  props.goal.description = newUpdates.description;
  props.goal.status = newUpdates.status;
  // do not send updates with empty strings, titles and descriptions cannot be empty
  if (!newUpdates.title || !newUpdates.description) return;

  // in the future we want to use a debouncer to reduce the number of api calls
  await updateGoalApi(props.goal.id, {
    title: newUpdates.title,
    description: newUpdates.description,
    status: newUpdates.status,
  });
});

const statuses = [
  {
    id: 1,
    name: "Not Complete",
    value: "not_complete",
  },
  {
    id: 2,
    name: "Complete",
    value: "complete",
  },
];

const statusMap = { not_complete: "Not Complete", complete: "Complete" };
</script>
<template>
  <header
    @click="openEditingDialog()"
    class="flex p-4 w-full h-full bg-gray-700 hover:cursor-pointer hover:bg-gray-600 gap-x-2 items-center rounded-sm"
    :class="{
      'bg-green-600 hover:bg-green-700': props.goal.status === 'complete',
    }"
  >
    <CheckOutline :onclick="(e: MouseEvent) => handleToggleStatus(e)" />
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
        @close="setIsEditing(false)"
      >
        <DialogPanel class="w-full sm:w-1/2">
          <Box
            gap="gap-4"
            shadow="shadow-lg"
            flex-direction="col"
            padding="p-8"
            height="h-full"
            width="w-full"
            class="border-white hover:cursor-default shadow-gray-400"
          >
            <InputField
              type="text"
              v-model="updates.title"
              containerWidth="w-full"
              class="text-3xl text-gray-300"
            >
              <template #right>
                <XMark
                  :onclick="(_) => setIsEditing(false)"
                  class="sm:ml-auto hover:cursor-pointer"
                />
              </template>
            </InputField>
            <Box flex-direction="col">
              <Text size="xl" as="p">Status</Text>
              <Box flex-direction="col" gap="gap-2">
                <Listbox v-model="updates.status">
                  <ListboxButton
                    :class="{
                      'w-56 h-8 text-center rounded-lg text-gray-600': true,
                      'bg-green-400': updates.status === 'complete',
                      'bg-orange-400': updates.status !== 'complete',
                    }"
                    ><Text color="dark">{{
                      statusMap[updates.status]
                    }}</Text></ListboxButton
                  >
                  <ListboxOptions class="">
                    <ListboxOption
                      v-for="status in statuses"
                      :key="status.id"
                      :value="status.value"
                      :disabled="status.value === updates.status"
                      :class="{
                        'w-56 hover:cursor-pointer h-8 bg-gray-300 text-gray-600 text-center hover:bg-gray-400 rounded-lg inline-flex items-center justify-center': true,
                        hidden: updates.status === status.value,
                      }"
                    >
                      <Text color="dark">{{ status.name }}</Text>
                    </ListboxOption>
                  </ListboxOptions>
                </Listbox>
              </Box>
            </Box>
            <Box flex-direction="col">
              <Text size="xl" as="p">Description</Text>
              <textarea
                v-model="updates.description"
                class="w-full bg-gray-300 focus:outline-none h-64 p-2 text-gray-600"
              />
            </Box>
            <Button
              variant="secondary"
              class="hover:bg-red-600 h-10"
              @click="setIsDeleting(true)"
            >
              <Text weight="semibold" size="base">Delete This Goal</Text>
            </Button>
            <DeleteModal
              :is-open="isDeleting"
              :set-opener="setIsDeleting"
              delete-message="Are you sure you want to delete this goal?"
              :delete-function="handleDeleteGoal"
            />
          </Box>
        </DialogPanel>
      </Dialog>
    </TransitionRoot>
  </section>
</template>
