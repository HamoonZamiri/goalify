<script setup lang="ts">
import { Dialog, DialogPanel } from "@headlessui/vue";

type Props = {
	isOpen: boolean;
	setOpener: (value: boolean) => void;
	deleteMessage: string;
	deleteFunction: (e: MouseEvent) => void;
};

const props = defineProps<Props>();
</script>
<template>
	<Dialog
		class="absolute inset-0 h-screen flex justify-center items-center hover:cursor-pointer z-20 w-screen bg-opacity-10 rounded-lg"
		:open="props.isOpen"
		@close="setOpener(false)"
	>
		<DialogPanel class="w-[25%] h-[25%]">
			<div class="border-2 border-white bg-gray-800 flex flex-col p-4 gap-y-4">
				<p class="text-xl text-gray-300 text-center">
					{{ props.deleteMessage }}
				</p>
				<div class="flex justify-center gap-4 text-gray-300">
					<button
						@click="setOpener(false)"
						class="bg-gray-400 text-gray-700 rounded-lg p-2 hover:bg-gray-500"
					>
						Cancel
					</button>
					<button
						class="bg-red-400 text-gray-700 rounded-lg p-2 hover:bg-red-500"
						@click="deleteFunction"
					>
						Yes, delete it
					</button>
				</div>
			</div>
		</DialogPanel>
	</Dialog>
</template>
