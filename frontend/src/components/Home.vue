<script setup lang="ts">
import GoalCategoryCard from "./goals/GoalCategoryCard.vue";
import { onMounted, ref } from "vue";
import CreateGoalCategoryDialog from "./goals/CreateGoalCategoryDialog.vue";
import type { GoalCategory } from "@/utils/schemas";
import { ApiClient } from "@/utils/api";
import ModalForm from "./ModalForm.vue";
import CreateGoalCategoryForm from "./goals/CreateGoalCategoryForm.vue";
import CreateCategoryButton from "./goals/CreateCategoryButton.vue";
import goalState from "@/state/goals";
const error = ref<string | null>(null);
const isLoading = ref<boolean>(true);

onMounted(async () => {
  const res = await ApiClient.getUserGoalCategories();
  if (typeof res === "string") {
    error.value = res;
    return;
  }

  goalState.categories = res.data;
  isLoading.value = false;
});
</script>

<template>
  <div v-if="isLoading">
    <v-icon name="co-reload" animation="spin" />
  </div>
  <div
    v-else
    class="flex flex-col items-center sm:items-start px-6 w-auto bg-slate-50"
  >
    <section class="flex-col sm:flex-row flex gap-4 w-auto">
      <div class="w-full sm:w-1/3" v-for="cat in goalState.categories">
        <GoalCategoryCard :goalCategory="cat" />
      </div>
      <ModalForm
        :FormComponent="CreateGoalCategoryForm"
        :OpenerComponent="CreateCategoryButton"
      />
    </section>
  </div>
</template>
