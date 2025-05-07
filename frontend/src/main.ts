import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import "./assets/index.css";
import Vue3Toastify, { type ToastContainerOptions } from "vue3-toastify";
import "vue3-toastify/dist/index.css";

const app = createApp(App);

app.use(router);

const toastifyOptions: ToastContainerOptions = {
	autoClose: 3000,
	position: "top-right",
	theme: "dark",
};
app.use(Vue3Toastify, toastifyOptions);

app.mount("#app");
