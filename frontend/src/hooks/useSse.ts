import { onUnmounted, ref } from "vue";

export function useSSE(url: string) {
  const eventSource = ref<EventSource | null>(null);
  const messages = ref<string[]>([]);

  const connect = () => {
    if (eventSource.value) {
      return;
    }

    const es = new EventSource(url);
    es.onopen = () => {
      console.log("connected");
    };
    es.onerror = () => {
      console.log("error");
    };
    es.onmessage = (event) => {
      console.log(event);
      console.log(event.data);
      messages.value.push(event.data);
    };

    eventSource.value = es;
  };

  onUnmounted(() => {
    if (eventSource.value) {
      eventSource.value.close();
    }
  });

  return {
    eventSource,
    connect,
    messages,
  };
}
