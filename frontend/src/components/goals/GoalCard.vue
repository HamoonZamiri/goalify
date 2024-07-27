<script setup lang="ts">
import type { Goal } from "@/utils/schemas";
import { ref, h, watch, reactive } from "vue";
import { Dialog, DialogPanel } from "@headlessui/vue";
import { ApiClient } from "@/utils/api";
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

function setIsEditing(value: boolean) {
  isEditing.value = value;
}

function getStatus(status: string) {
  if (status === "complete") {
    return h("p", { class: "rounded-lg bg-green-400 px-2" }, "Complete");
  }
  return h("p", { class: "rounded-lg bg-orange-400 px-2" }, "Not Complete");
}

// watches for updates on the goal title and description
watch(updates, async (newUpdates) => {
  // syncronhize state passed in from props with local reactive updates
  props.goal.title = newUpdates.title;
  props.goal.description = newUpdates.description;
  // do not send updates with empty strings, titles and descriptions cannot be empty
  if (!newUpdates.title || !newUpdates.description) return;

  // in the future we want to use a debouncer to reduce the number of api calls
  await ApiClient.updateGoal(props.goal.id, {
    title: newUpdates.title,
    description: newUpdates.description,
  });
});
</script>
<template>
  <header
    @click="setIsEditing(true)"
    class="flex p-4 w-full h-full bg-gray-700 hover:cursor-pointer hover:bg-gray-600 gap-x-2 items-center rounded-sm"
  >
    <v-icon
      class="hover:cursor-pointer"
      name="io-checkmark-circle-outline"
      animation="wrench"
      hover
    />
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
            <component :is="getStatus(props.goal.status)" />
          </div>
          <textarea
            v-model="updates.description"
            class="w-full bg-gray-300 focus:outline-none h-64 p-2"
          />
        </div>
      </DialogPanel>
    </Dialog>
  </section>
</template>
