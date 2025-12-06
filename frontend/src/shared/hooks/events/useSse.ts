import { useQueryClient } from "@tanstack/vue-query";
import { readonly, ref } from "vue";
import { toast } from "vue3-toastify";
import { z } from "zod";
import { categoryKeys } from "@/features/goals/queries";
import { API_BASE, events } from "@/utils/constants";
import useAuth from "../auth/useAuth";

const xpUpdateSchema = z.object({
	xp: z.number(),
	level_id: z.number(),
});

const eventSource = ref<EventSource>();
const MAX_RECONNECT_ATTEMPTS = 3;
const reconnectAttempts = ref(0);

/**
 * Hook for managing Server-Sent Events connection
 */
export function useSSE() {
	const queryClient = useQueryClient();
	const { getUser, setUser } = useAuth();

	/**
	 * Closes the current SSE connection
	 */
	const closeConnection = () => {
		if (eventSource.value) {
			eventSource.value.close();
			eventSource.value = undefined;
			reconnectAttempts.value = 0;
		}
	};

	/**
	 * Sets up all event listeners on an EventSource instance
	 */
	const setupEventListeners = (es: EventSource) => {
		es.onerror = () => {
			closeConnection();

			if (reconnectAttempts.value < MAX_RECONNECT_ATTEMPTS) {
				reconnectAttempts.value++;
				setTimeout(() => {
					connect(`${API_BASE}/events?token=${getUser()?.access_token}`);
				}, 1000 * reconnectAttempts.value);
				return;
			}
			toast.error("Failed to connect to the server. Please refresh the page.");
		};

		es.addEventListener(events.DEFAULT_GOAL_CREATED, () => {
			queryClient.invalidateQueries({ queryKey: categoryKeys.all });
		});

		es.addEventListener(events.XP_UPDATED, (event) => {
			const json = JSON.parse(event.data);
			const parsedData = xpUpdateSchema.parse(json);
			const user = getUser();
			if (user) {
				setUser({ ...user, ...parsedData });
			}
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
