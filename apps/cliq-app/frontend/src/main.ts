import { createApp } from "vue";
import App from "@/App.vue";
import "@/styles/global.css";

// PrimeVue
import PrimeVue from "primevue/config";
import Aura from "@primeuix/themes/aura";
import "primeicons/primeicons.css";

import PrimeComponents from "@/components/PrimeComponents";

const app = createApp(App);
app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: {
      darkModeSelector: ".disable-app-dark",
    },
  },
});

app.use(PrimeComponents);
app.mount("#app");
