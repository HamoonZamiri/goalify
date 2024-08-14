import { onUnmounted, ref } from "vue";

export function useSSE(url: string) {
  const eventSource = ref<EventSource | null>(null);

  const connect = () => {
    if (eventSource.value) {
      return;
    }

    eventSource.value = new EventSource(url);
    eventSource.value.onopen = () => {
      console.log("connected");
    };
    eventSource.value.onerror = () => {
      console.log("error");
    };
    eventSource.value.onmessage = (event) => {
      console.log(event);
      console.log(event.data);
    };
  };

  onUnmounted(() => {
    if (eventSource.value) {
      eventSource.value.close();
    }
  });

  return {
    eventSource,
    connect,
  };
}
