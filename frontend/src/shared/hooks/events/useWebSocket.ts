import { onUnmounted, ref } from "vue";
import { z } from "zod";
import { useQueryClient } from "@tanstack/vue-query";
import { categoryKeys } from "@/features/goals/queries";
import { GoalSchema, GoalCategorySchema } from "@/features/goals/schemas";
import { UserSchema } from "@/features/auth/schemas";
import { events } from "@/utils/constants";

function createEventSchema<TData extends z.ZodTypeAny>(schema: TData) {
	return z.object({
		event_type: z.string(),
		data: schema,
		user_id: z.string().uuid(),
	});
}

const UserEventSchema = createEventSchema(UserSchema);
const GoalEventSchema = createEventSchema(GoalSchema);
const GoalCategoryEventSchema = createEventSchema(GoalCategorySchema);

export const EventSchemas = {
	[events.USER_CREATED]: UserEventSchema,
	[events.GOAL_CREATED]: GoalEventSchema,
	[events.GOAL_CATEGORY_CREATED]: GoalCategoryEventSchema,
	[events.USER_UPDATED]: UserEventSchema,
	[events.DEFAULT_GOAL_CREATED]: GoalEventSchema,
} as const;

export default function useWebSocket(url: string) {
	const queryClient = useQueryClient();
	const websocket = ref<WebSocket | null>(null);

	function handleDefaultGoalCreated(
		event: z.infer<(typeof EventSchemas)[typeof events.DEFAULT_GOAL_CREATED]>,
	) {
		queryClient.invalidateQueries({ queryKey: categoryKeys.all });
	}

	function handleEvent(event: MessageEvent) {
		const json = JSON.parse(event.data);
		const eventType = json.event_type as string;
		switch (eventType) {
			case events.DEFAULT_GOAL_CREATED: {
				const parsedEvent =
					EventSchemas[events.DEFAULT_GOAL_CREATED].parse(json);
				handleDefaultGoalCreated(parsedEvent);
				break;
			}
			default:
				console.log("unhandled event:", eventType);
		}
	}

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
