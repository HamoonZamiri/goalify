import { events } from "@/utils/constants";
import { Schemas, type User } from "@/utils/schemas";
import { readonly, ref } from "vue";
import useGoals from "@/hooks/goals/useGoals";
import useAuth from "../auth/useAuth";
import { z } from "zod";
import { toast } from "vue3-toastify";

const xpUpdateSchema = z.object({
	xp: z.number(),
	level_id: z.number(),
});

const eventSource = ref<EventSource>();

/**
 * Hook for managing Server-Sent Events connection
 */
export function useSSE() {
	const { addGoal } = useGoals();
	const { getUser, setUser } = useAuth();

	/**
	 * Closes the current SSE connection
	 */
	const closeConnection = () => {
		if (eventSource.value) {
			eventSource.value.close();
			eventSource.value = undefined;
		}
	};

	/**
	 * Sets up all event listeners on an EventSource instance
	 */
	const setupEventListeners = (es: EventSource) => {
		es.onerror = () => {
			toast.warning("There was an issue connecting to the server!");
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

		es.addEventListener(events.XP_UPDATED, (event) => {
			const json = JSON.parse(event.data);
			const parsedData = xpUpdateSchema.parse(json);
			const user = getUser() as User;
			setUser({ ...user, ...parsedData });
		});
	};

	/**
	 * Establishes SSE connection to the given URL
	 */
	const connect = (url: string) => {
		if (eventSource.value) {
			return;
		}

		const es = new EventSource(url);
		setupEventListeners(es);
		eventSource.value = es;
	};

	/**
	 * Reconnects SSE with a new URL (e.g., after token refresh)
	 */
	const reconnect = (url: string) => {
		closeConnection();
		connect(url);
	};

	return {
		eventSource: readonly(eventSource),
		connect,
		reconnect,
		closeConnection,
	};
}
