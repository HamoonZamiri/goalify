<script setup lang="ts">
import { type InputTypeHTMLAttribute, ref } from "vue";
import {
  type Height,
  type TextColor,
  textColorMap,
  type Width,
} from "@/utils/tailwind";
import Box from "./Box.vue";

defineOptions({
  inheritAttrs: false,
});

type InputFieldProps = {
  value?: string | number;
  accept?: string;
  alt?: string;
  autocomplete?: string;
  disabled?: boolean;
  name?: string;
  placeholder?: string;
  type?: InputTypeHTMLAttribute;
  containerWidth?: Width;
  width?: Width;
  height?: Height;
  textColor?: TextColor;
  class?: string;
  bg?: "transparent" | "primary";
  errorslot?: boolean;
  compact?: boolean;
  as?: "input" | "textarea";
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
  as: "input",
});

const emit = defineEmits<{
  input: [value: string | number | null];
  blur: [event: Event];
}>();

const baseClass = props.compact
  ? "focus:outline-none placeholder-gray-400 rounded border-0 text-center min-w-0"
  : "focus:outline-none placeholder-gray-400 sm:flex-1 rounded-lg border-0";

const inputRef = ref<HTMLInputElement | HTMLTextAreaElement | null>(null);

function handleBeforeInput(e: Event) {
  if (props.type !== "number") return;

  const inputEvent = e as InputEvent;
  const data = inputEvent.data;
  if (!data) return; // Allow deletions, etc.

  // Only allow digits, decimal point, and minus sign
  const isValid = /^[\d.-]$/.test(data);

  if (!isValid) {
    e.preventDefault();
  }
}

function handleInput(e: Event) {
  const target = e.target as HTMLInputElement | HTMLTextAreaElement;
  const raw = target.value;

  if (props.type === "number") {
    if (raw === "" || raw === "-") {
      emit("input", null);
      return;
    }

    const parsed = Number(raw);
    if (Number.isNaN(parsed)) {
      emit("input", null);
      return;
    }

    emit("input", parsed);
  } else {
    emit("input", raw);
  }
}

function handleBlur(e: Event) {
  emit("blur", e);
}

const sharedClasses = [
  baseClass,
  bgClasses[props.bg],
  textColorMap[props.textColor],
  props.bg !== "transparent" ? "px-1.5 py-1.5" : "",
  props.compact ? "w-12 text-xs" : "",
  props.class,
];
</script>

<template>
  <Box class="bg-inherit" flex-direction="col" :width="props.containerWidth">
    <slot name="label" />
    <Box flex-direction="row" class="gap-1">
      <slot name="left" />
      <input
        v-if="props.as === 'input'"
        ref="inputRef"
        :class="sharedClasses"
        :id="props.name"
        :name="props.name"
        :type="props.type === 'number' ? 'text' : props.type"
        :inputmode="props.type === 'number' ? 'numeric' : undefined"
        :value="props.value ?? ''"
        :placeholder="props.placeholder"
        :disabled="props.disabled"
        :accept="props.accept"
        :autocomplete="props.autocomplete"
        @beforeinput="handleBeforeInput"
        @input="handleInput"
        @blur="handleBlur"
      />

      <textarea
        v-else-if="props.as === 'textarea'"
        ref="inputRef"
        :class="sharedClasses"
        :id="props.name"
        :name="props.name"
        :placeholder="props.placeholder"
        :disabled="props.disabled"
        :value="props.value ?? ''"
        @input="handleInput"
        @blur="handleBlur"
      />
      <slot name="right" />
    </Box>
    <Box v-if="props.errorslot" class="min-h-6 bg-inherit">
      <slot name="error" />
    </Box>
  </Box>
</template>
