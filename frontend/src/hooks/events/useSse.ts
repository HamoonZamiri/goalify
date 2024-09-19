import { events } from "@/utils/constants";
import { Schemas, type User } from "@/utils/schemas";
import { onUnmounted, ref } from "vue";
import useGoals from "@/hooks/goals/useGoals";
import useAuth from "../auth/useAuth";
import { z } from "zod";

const xpUpdateSchema = z.object({
  xp: z.number(),
  level_id: z.number(),
});

export function useSSE(url: string) {
  const { addGoal } = useGoals();
  const { getUser, setUser } = useAuth();
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

    es.addEventListener(events.SSE_CONNECTED, () => {
      console.log("initial sse event");
    });

    es.addEventListener(events.XP_UPDATED, (event) => {
      const json = JSON.parse(event.data);
      const parsedData = xpUpdateSchema.parse(json);
      const user = getUser() as User;
      setUser({ ...user, ...parsedData });
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
