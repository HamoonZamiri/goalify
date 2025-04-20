import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import "./assets/index.css";
import { OhVueIcon, addIcons } from "oh-vue-icons";
import {
  IoCheckmarkCircleOutline,
  PxNotesPlus,
  CoPlus,
  CoReload,
} from "oh-vue-icons/icons";
import Vue3Toastify, { type ToastContainerOptions } from "vue3-toastify";
import "vue3-toastify/dist/index.css";

addIcons(IoCheckmarkCircleOutline, PxNotesPlus, CoPlus, CoReload);
const app = createApp(App);

app.use(router);

const toastifyOptions: ToastContainerOptions = {
  autoClose: 3000,
  position: "top-right",
  theme: "dark",
};
app.use(Vue3Toastify, toastifyOptions);

app.component("v-icon", OhVueIcon);
app.mount("#app");
