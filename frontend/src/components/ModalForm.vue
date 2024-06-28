<script setup lang="ts">
import { ref } from "vue";
import { Dialog, DialogPanel } from "@headlessui/vue";

const ModalFormProps = defineProps({
  FormComponent: Object,
  OpenerComponent: Object,
});
const isOpen = ref(false);

function setIsOpen(value: boolean) {
  isOpen.value = value;
}
</script>

<template>
  <div
    @click="setIsOpen(true)"
    class="hover:cursor-pointer hover:text-blue-500"
  >
    <component
      :is-open="isOpen"
      :set-is-open="setIsOpen"
      :is="ModalFormProps.OpenerComponent"
    />
  </div>
  <Dialog
    className="absolute inset-0 h-screen flex justify-center items-center
    hover:cursor-pointer bg-gray-600 z-10 w-screen bg-opacity-90 rounded-lg p-4"
    :open="isOpen"
    @close="setIsOpen(false)"
  >
    <DialogPanel>
      <component
        :is-open="isOpen"
        :set-is-open="setIsOpen"
        :is="ModalFormProps.FormComponent"
      />
    </DialogPanel>
  </Dialog>
</template>
