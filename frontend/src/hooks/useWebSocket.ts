import { onUnmounted, ref } from "vue";

export default function useWebSocket(url: string) {
  const websocket = ref<WebSocket | null>(null);

  const connect = () => {
    if (websocket.value) {
      return;
    }
    const ws = new WebSocket(url);
    ws.onopen = () => {
      console.log("connected");
    };
    ws.onerror = () => {
      console.log("error");
    };
    ws.onmessage = (event) => {
      console.log(event);
      console.log(event.data);
    };
    websocket.value = ws;
  };

  onUnmounted(() => {
    if (websocket.value) {
      websocket.value.close();
    }
  });

  return { connect };
}
