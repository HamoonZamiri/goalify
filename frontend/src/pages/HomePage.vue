<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import { useGoalCategories } from "@/features/goals/queries";
import {
	GoalCategoryCard,
	CreateCategoryButton,
} from "@/features/goals/components";
import { CreateGoalCategoryForm } from "@/features/goals/forms";
import useAuth from "@/shared/hooks/auth/useAuth";
import { useSSE } from "@/shared/hooks/events/useSse";
import { API_BASE } from "@/utils/constants";
import { Box } from "@/shared/components/ui";
import { ArrowPath } from "@/shared/components/icons";
import { ModalForm } from "@/shared/components/modals";

const isCreateCategoryDialogOpen = ref(false);

const { getUser } = useAuth();
const { data: categories, isLoading, error } = useGoalCategories();

const { connect, closeConnection } = useSSE();

onMounted(() => {
	connect(`${API_BASE}/events?token=${getUser()?.access_token}`);
});

onUnmounted(() => {
	closeConnection();
});
</script>

<template>
  <ArrowPath class="animate-spin" v-if="isLoading" />
  <Box
    v-else-if="error"
    height="h-full"
    bg="darkest"
    class="items-center justify-center"
  >
    <p class="text-red-500">Error loading categories: {{ error.message }}</p>
  </Box>
  <Box
    v-else
    height="h-full"
    bg="darkest"
    class="items-center sm:items-start px-6 w-full overflow-hidden"
  >
    <Box
      gap="gap-4"
      bg="darkest"
      width="w-full"
      class="flex-grow sm:flex-row overflow-x-auto"
    >
      <Box
        width="w-full"
        bg="darkest"
        class="sm:w-1/2 lg:w-1/3 flex-shrink-0"
        v-for="cat in categories"
        :key="cat.id"
      >
        <GoalCategoryCard :goalCategory="cat" />
      </Box>
      <Box bg="darkest" flex-direction="row">
        <CreateCategoryButton
          class="hover:cursor-pointer"
          @click="isCreateCategoryDialogOpen = true"
        />
        <ModalForm
          v-model="isCreateCategoryDialogOpen"
          @close="isCreateCategoryDialogOpen = false"
        >
          <CreateGoalCategoryForm @close="isCreateCategoryDialogOpen = false" />
        </ModalForm>
      </Box>
    </Box>
  </Box>
</template>
