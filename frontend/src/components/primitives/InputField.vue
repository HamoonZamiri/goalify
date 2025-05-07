<script setup lang="ts">
import type { InputTypeHTMLAttribute } from "vue";
import Box from "./Box.vue";
import {
  type Width,
  type Height,
  type TextColor,
  textColorMap,
} from "@/utils/tailwind";

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
  containerWidth?: Width;
  width?: Width;
  height?: Height;
  textColor?: TextColor;
  class?: string;
  bg?: "transparent" | "primary";
  errorslot?: boolean;
  compact?: boolean;
};

const bgClasses = {
  transparent: "bg-transparent",
  primary: "bg-gray-300",
};

const props = withDefaults(defineProps<InputFieldProps>(), {
  type: "text",
  bg: "transparent",
  textColor: "light",
  height: "h-10",
  width: "w-full",
  compact: false,
});

const emit = defineEmits(["update:modelValue"]);
const baseClass = props.compact
  ? "focus:outline-none placeholder-gray-400 rounded border-0 text-center min-w-0"
  : "focus:outline-none placeholder-gray-400 sm:flex-1 rounded-lg border-0";
</script>

<template>
  <Box class="bg-inherit" flex-direction="col" :width="props.containerWidth">
    <slot name="label" />
    <Box flex-direction="row" class="gap-1">
      <slot name="left" />
      <input
        :class="[
          baseClass,
          bgClasses[props.bg],
          textColorMap[props.textColor],
          props.bg !== 'transparent' ? 'px-1.5 py-1.5' : '',
          props.compact ? 'w-12 text-xs' : '',
          props.class,
        ]"
        v-bind="props"
        :value="props.modelValue"
        @input="
          (e) => {
            const raw = (e.target as HTMLInputElement).value;
            const parsed =
              props.type === 'number' ? (raw === '' ? null : Number(raw)) : raw;
            emit('update:modelValue', parsed as any);
          }
        "
      />
      <slot name="right" />
    </Box>
    <Box v-if="props.errorslot" class="min-h-6 bg-inherit">
      <slot name="error" />
    </Box>
  </Box>
</template>
