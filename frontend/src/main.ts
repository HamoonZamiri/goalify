import { VueQueryPlugin } from "@tanstack/vue-query";
import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import "./assets/index.css";
import Vue3Toastify, { type ToastContainerOptions } from "vue3-toastify";
import "vue3-toastify/dist/index.css";
import { queryClient } from "@/shared/api/query-client";

const app = createApp(App);

app.use(router);
app.use(VueQueryPlugin, { queryClient });

const toastifyOptions: ToastContainerOptions = {
	autoClose: 3000,
	position: "top-right",
	theme: "dark",
};
app.use(Vue3Toastify, toastifyOptions);

app.mount("#app");
