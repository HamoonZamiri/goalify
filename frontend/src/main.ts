import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import "./assets/index.css";
import { OhVueIcon, addIcons } from "oh-vue-icons";
import {
  IoCheckmarkCircleOutline,
  PxNotesPlus,
  CoPlus,
} from "oh-vue-icons/icons";
addIcons(IoCheckmarkCircleOutline, PxNotesPlus, CoPlus);
const app = createApp(App);

app.use(router);

app.component("v-icon", OhVueIcon);
app.mount("#app");
