import goalState from "@/state/goals";
import { events } from "@/utils/constants";
import { Schemas } from "@/utils/schemas";
import { onUnmounted, ref } from "vue";
import { z } from "zod";

function createEventSchema<TData extends z.ZodTypeAny>(schema: TData) {
  return z.object({
    event_type: z.string(),
    data: schema,
    user_id: z.string().uuid(),
  });
}

const UserEventSchema = createEventSchema(Schemas.UserSchema);
const GoalEventSchema = createEventSchema(Schemas.GoalSchema);
const GoalCategoryEventSchema = createEventSchema(Schemas.GoalCategorySchema);

export const EventSchemas = {
  [events.USER_CREATED]: UserEventSchema,
  [events.GOAL_CREATED]: GoalEventSchema,
  [events.GOAL_CATEGORY_CREATED]: GoalCategoryEventSchema,
  [events.USER_UPDATED]: UserEventSchema,
  [events.DEFAULT_GOAL_CREATED]: GoalEventSchema,
} as const;

export function handleDefaultGoalCreated(
  event: z.infer<(typeof EventSchemas)[typeof events.DEFAULT_GOAL_CREATED]>,
) {
  goalState.addGoal(event.data.category_id, event.data);
}
export default function useWebSocket(url: string) {
  const websocket = ref<WebSocket | null>(null);

  const connect = () => {
    if (websocket.value) {
      return;
    }
    const ws = new WebSocket(url);
    ws.onopen = () => {
      console.log("connected to websocket server for event processing");
    };
    ws.onerror = () => {
      console.log("error");
    };
    ws.onmessage = (event) => {
      handleEvent(event);
    };
    websocket.value = ws;
  };

  onUnmounted(() => {
    if (websocket.value) {
      websocket.value.close();
      websocket.value = null;
    }
  });

  return { connect };
}

function handleEvent(event: MessageEvent) {
  const json = JSON.parse(event.data);
  const eventType = json.event_type as string;
  switch (eventType) {
    case events.DEFAULT_GOAL_CREATED: {
      const parsedEvent = EventSchemas[events.DEFAULT_GOAL_CREATED].parse(json);
      handleDefaultGoalCreated(parsedEvent);
      break;
    }
    default:
      console.log("unhandled event:", eventType);
  }
}
