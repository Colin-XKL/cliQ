import { App } from "vue";
import Toast from "primevue/toast";
import ToastService from "primevue/toastservice";

export default {
  install: (app: App) => {
    app.use(ToastService);
    app.component("Toast", Toast);
  },
};
