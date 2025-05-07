<script setup lang="ts">
import type { InputHTMLAttributes, InputTypeHTMLAttribute } from "vue";
import Box from "./Box.vue";
import type { Width } from "@/utils/tailwind";

defineOptions({
  inheritAttrs: false,
});

type InputFieldProps = {
  modelValue?: string | number;
  accept?: string;
  alt?: string;
  autocomplete?: string;
  disabled?: boolean;
  name?: string;
  placeholder?: string;
  type?: InputTypeHTMLAttribute;
  value?: unknown;
  width?: Width;
  class?: string;
};

const props = defineProps<InputFieldProps>();
const emit = defineEmits(["update:modelValue"]);
const baseClass =
  "bg-transparent focus:outline-none text-gray-200 placeholder-gray-400";
</script>

<template>
  <Box flex-direction="row">
    <slot name="left" />
    <input
      :class="[props.class, baseClass]"
      v-bind="props"
      :value="props.modelValue"
      @input="
        (e) => emit('update:modelValue', (e.target as HTMLInputElement).value)
      "
    />
    <slot name="right" />
  </Box>
</template>
