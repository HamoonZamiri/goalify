<script setup lang="ts">
import { ref } from "vue";
import {
  Dialog,
  DialogDescription,
  DialogPanel,
  DialogTitle,
} from "@headlessui/vue";
type CreateCategoryForm = {
  title: string;
  xp_per_goal: number;
};

const formData = ref<CreateCategoryForm>({
  title: "",
  xp_per_goal: 0,
});

const isOpen = ref(false);

function setIsOpen(value: boolean) {
  isOpen.value = value;
}

function handleSubmit(e: MouseEvent) {
  e.preventDefault();
  setIsOpen(false);
}
</script>

<template>
  <div
    @click="setIsOpen(true)"
    class="hover:cursor-pointer hover:text-blue-500"
  >
    <v-icon name="co-plus" />
    <span class="text-xl text-blue-400 font-semibold" @click="setIsOpen(true)">
      Add Category
    </span>
  </div>
  <Dialog
    className="absolute inset-0 h-screen flex justify-center items-center hover:cursor-pointer bg-gray-600 z-10 w-screen  bg-opacity-90 rounded-lg p-4"
    :open="isOpen"
    @close="setIsOpen(false)"
  >
    <DialogPanel>
      <form
        class="rounded-lg border bg-white p-3 w-[400px] grid grid-cols-1 gap-8 hover:cursor-default"
      >
        <DialogTitle class="font-semibold"
          >Create a New Goal/Task Category</DialogTitle
        >
        <div class="grid grid-cols-1 gap-4">
          <label>Title:</label>
          <input
            type="text"
            v-model="formData.title"
            class="block w-full rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-200 sm:text-sm sm:leading-6"
          />
        </div>
        <div>
          <label>XP/goal:</label>
          <input
            type="number"
            min="1"
            max="100"
            class="block w-full rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-200 sm:text-sm sm:leading-6"
          />
        </div>
        <button
          @click="handleSubmit"
          class="block w-full hover:bg-blue-200 bg-blue-100 rounded-lg h-10"
        >
          Add Category
        </button>
      </form>
    </DialogPanel>
  </Dialog>
</template>
