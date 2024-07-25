<script setup lang="ts">
import { ref } from "vue";
import { Dialog, DialogPanel } from "@headlessui/vue";

const ModalFormProps = defineProps({
  FormComponent: Object,
  OpenerComponent: Object,
  formProps: Object,
});
const isOpen = ref(false);

function setIsOpen(value: boolean) {
  isOpen.value = value;
}
</script>

<template>
  <div class="hover:cursor-pointer" @click="setIsOpen(true)">
    <component
      :is-open="isOpen"
      :set-is-open="setIsOpen"
      :is="ModalFormProps.OpenerComponent"
    />
  </div>
  <Dialog
    className="absolute inset-0 h-screen flex justify-center items-center
    hover:cursor-pointer bg-gray-600 z-10 w-screen bg-opacity-20 rounded-lg p-4"
    :open="isOpen"
    @close="setIsOpen(false)"
  >
    <DialogPanel>
      <component
        :props="ModalFormProps.formProps"
        :is-open="isOpen"
        :set-is-open="setIsOpen"
        :is="ModalFormProps.FormComponent"
      />
    </DialogPanel>
  </Dialog>
</template>
