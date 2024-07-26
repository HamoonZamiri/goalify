<script setup lang="ts">
import type { Goal } from "@/utils/schemas";
import { ref, h } from "vue";
import { Dialog, DialogPanel } from "@headlessui/vue";
const props = defineProps<{
  goal: Goal;
}>();

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
            v-model="props.goal.title"
            class="w-full bg-gray-800 text-gray-200 focus:outline-none text-3xl"
          />
          <div class="flex gap-x-24 w-full text-gray-200">
            <p class="text-xl">Status</p>
            <component :is="getStatus(props.goal.status)" />
          </div>
          <textarea
            v-model="props.goal.description"
            class="w-full bg-gray-300 focus:outline-none h-64 p-2"
          />
        </div>
      </DialogPanel>
    </Dialog>
  </section>
</template>
