<script setup lang="ts">
import GoalCategoryCard from "@/components/goals/cards/GoalCategoryCard.vue";
import { computed, onMounted, ref, watch } from "vue";
import { type ErrorResponse, ApiClient } from "@/utils/api";
import ModalForm from "@/components/ModalForm.vue";
import CreateGoalCategoryForm from "@/components/goals/forms/CreateGoalCategoryForm.vue";
import CreateCategoryButton from "@/components/goals/buttons/CreateCategoryButton.vue";
import authState from "@/state/auth";
import { useSSE } from "@/hooks/events/useSse";
import useGoals from "@/hooks/goals/useGoals";

// State
const error = ref<ErrorResponse | null>(null);
const isLoading = ref<boolean>(true);
const { connect } = useSSE(
  `http://localhost:8080/api/events?token=${authState.getUser()?.access_token}`,
);
const { setCategories, categoryState } = useGoals();

onMounted(async () => {
  const res = await ApiClient.getUserGoalCategories();
  if (ApiClient.isError(res)) {
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
