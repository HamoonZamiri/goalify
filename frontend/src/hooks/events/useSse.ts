import { events } from "@/utils/constants";
import { Schemas, type User } from "@/utils/schemas";
import { onUnmounted, ref } from "vue";
import useGoals from "@/hooks/goals/useGoals";
import useAuth from "../auth/useAuth";
import { z } from "zod";
import { toast } from "vue3-toastify";

const xpUpdateSchema = z.object({
	xp: z.number(),
	level_id: z.number(),
});

// moving the event source outside the hook to make it a global singleton
const eventSource = ref<EventSource>();
export function useSSE(url: string) {
	const { addGoal } = useGoals();
	const { getUser, setUser } = useAuth();

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
			toast.warning("There was an issue connecting to the server!");

			console.error("error", event);
			closeConnection();
		};

		es.addEventListener(events.DEFAULT_GOAL_CREATED, (event) => {
			const json = JSON.parse(event.data);
			const parsedData = Schemas.GoalSchema.parse(json);
			toast.success(
				"Goal Category was created! We created a default example goal for you. You can delete it later!",
			);
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

	const closeConnection = () => {
		if (eventSource.value) {
			eventSource.value.close();
		}
		eventSource.value = undefined;
	};

	onUnmounted(() => {
		closeConnection();
	});

	return {
		eventSource,
		connect,
		closeConnection,
	};
}
