<script setup lang="ts">
import GoalCategoryCard from "@/components/goals/cards/GoalCategoryCard.vue";
import { onMounted, ref } from "vue";
import ModalForm from "@/components/ModalForm.vue";
import CreateGoalCategoryForm from "@/components/goals/forms/CreateGoalCategoryForm.vue";
import CreateCategoryButton from "@/components/goals/buttons/CreateCategoryButton.vue";
import { useSSE } from "@/hooks/events/useSse";
import useGoals from "@/hooks/goals/useGoals";
import useAuth from "@/hooks/auth/useAuth";
import type { ErrorResponse } from "@/utils/schemas";
import useApi from "@/hooks/api/useApi";
import { API_BASE } from "@/utils/constants";

// State
const { getUser } = useAuth();
const { getUserGoalCategories, isError } = useApi();
const error = ref<ErrorResponse | null>(null);
const isLoading = ref<boolean>(true);
const { connect } = useSSE(
  `${API_BASE}/events?token=${getUser()?.access_token}`,
);
const { setCategories, categoryState } = useGoals();

onMounted(async () => {
  const res = await getUserGoalCategories();
  if (isError(res)) {
    // in this case we are only expecting a message and not input validation errors
    error.value = res;
    return;
  }

  setCategories(res.data);
  isLoading.value = false;
  connect();
});
</script>

<template>
  <div v-if="isLoading">
    <v-icon name="co-reload" animation="spin" />
  </div>
  <div
    v-else
    class="flex flex-col h-full bg-gray-900 items-center sm:items-start px-6 w-full overflow-hidden"
  >
    <section
      class="flex-col flex-grow sm:flex-row flex gap-4 w-full overflow-x-auto"
    >
      <div
        class="w-full sm:w-1/3 flex-shrink-0"
        v-for="cat in categoryState.categories"
        key="cat.id"
      >
        <GoalCategoryCard :goalCategory="cat" />
      </div>
      <ModalForm
        :FormComponent="CreateGoalCategoryForm"
        :OpenerComponent="CreateCategoryButton"
      />
    </section>
  </div>
</template>
