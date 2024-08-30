import { events } from "@/utils/constants";
import { Schemas } from "@/utils/schemas";
import { onUnmounted, ref } from "vue";
import useGoals from "@/hooks/goals/useGoals";

export function useSSE(url: string) {
  const { addGoal } = useGoals();
  const eventSource = ref<EventSource | null>(null);

  const connect = () => {
    if (eventSource.value) {
      return;
    }

    const es = new EventSource(url);
    es.onopen = () => {
      console.log("connected");
      console.log("readystate:", es.readyState);
    };
    es.onerror = (event) => {
      console.error("error", event);
    };

    es.addEventListener(events.DEFAULT_GOAL_CREATED, (event) => {
      const json = JSON.parse(event.data);
      const parsedData = Schemas.GoalSchema.parse(json);
      addGoal(parsedData.category_id, parsedData);
    });

    eventSource.value = es;
  };

  onUnmounted(() => {
    if (eventSource.value) {
      eventSource.value.close();
    }
    eventSource.value = null;
  });

  return {
    eventSource,
    connect,
  };
}
