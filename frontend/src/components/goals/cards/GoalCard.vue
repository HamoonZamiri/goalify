<script setup lang="ts">
import type { Goal } from "@/utils/schemas";
import { ref, h, watch, reactive } from "vue";
import {
  Dialog,
  DialogPanel,
  Listbox,
  ListboxButton,
  ListboxOption,
  ListboxOptions,
} from "@headlessui/vue";
import useGoals from "@/hooks/goals/useGoals";
import useApi from "@/hooks/api/useApi";
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

async function handleDeleteGoal(e: MouseEvent) {
  e.preventDefault();

  await deleteGoalApi(props.goal.id);

  // remove goal from state
  deleteGoal(props.goal.category_id, props.goal.id);

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
    @click="setIsEditing(true)"
    class="flex p-4 w-full h-full bg-gray-700 hover:cursor-pointer hover:bg-gray-600 gap-x-2 items-center rounded-sm"
    :class="{
      'bg-green-600 hover:bg-green-700': props.goal.status === 'complete',
    }"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      stroke-width="1.5"
      stroke="currentColor"
      class="size-6"
      @click.stop="handleToggleStatus"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
      />
    </svg>
    <span class="font-semibold text-gray-300">{{ props.goal.title }}</span>
  </header>
  <section>
    <Dialog
      class="absolute inset-0 h-screen flex justify-end hover:cursor-pointer z-10 w-screen bg-opacity-10 rounded-lg"
      :open="isEditing"
      @close="setIsEditing(false)"
    >
      <DialogPanel class="w-1/2">
        <div
          class="flex flex-col gap-4 h-full p-8 border-white bg-gray-800 hover:cursor-default shadow-md shadow-gray-400"
        >
          <input
            v-model="updates.title"
            class="w-full bg-gray-800 text-gray-200 focus:outline-none text-3xl"
          />
          <div class="flex gap-x-24 w-full text-gray-200">
            <p class="text-xl">Status</p>
            <div class="">
              <Listbox v-model="updates.status">
                <ListboxButton
                  :class="{
                    'w-56 h-8 text-center rounded-lg text-gray-600': true,
                    'bg-green-400': updates.status === 'complete',
                    'bg-orange-400': updates.status !== 'complete',
                  }"
                  >{{ statusMap[updates.status] }}</ListboxButton
                >
                <ListboxOptions class="mt-1">
                  <ListboxOption
                    v-for="status in statuses"
                    :key="status.id"
                    :value="status.value"
                    :disabled="status.value === updates.status"
                    :class="{
                      'w-56 hover:cursor-pointer h-8 bg-gray-300 text-gray-600 text-center hover:bg-gray-400': true,
                      hidden: updates.status === status.value,
                    }"
                  >
                    {{ status.name }}
                  </ListboxOption>
                </ListboxOptions>
              </Listbox>
            </div>
          </div>
          <textarea
            v-model="updates.description"
            class="w-full bg-gray-300 focus:outline-none h-64 p-2"
          />
          <button
            class="w-full bg-red-400 hover:bg-red-500 text-gray-300 rounded-sm h-10"
            @click="setIsDeleting(true)"
          >
            Delete This Goal
          </button>

          <Dialog
            class="absolute inset-0 h-screen flex justify-center items-center hover:cursor-pointer z-20 w-screen bg-opacity-10 rounded-lg"
            :open="isDeleting"
            @close="setIsDeleting(false)"
          >
            <DialogPanel class="w-[25%] h-[25%]">
              <div
                class="border-2 border-white bg-gray-800 flex flex-col p-4 gap-y-4"
              >
                <p class="text-xl text-gray-300 text-center">
                  Are you sure you want to delete this goal?
                </p>
                <div class="flex justify-center gap-4 text-gray-300">
                  <button
                    @click="setIsDeleting(false)"
                    class="bg-gray-400 text-gray-700 rounded-lg p-2 hover:bg-gray-500"
                  >
                    Cancel
                  </button>
                  <button
                    class="bg-red-400 text-gray-700 rounded-lg p-2 hover:bg-red-500"
                    @click="handleDeleteGoal"
                  >
                    Yes, delete it
                  </button>
                </div>
              </div>
            </DialogPanel>
          </Dialog>
        </div>
      </DialogPanel>
    </Dialog>
  </section>
</template>
