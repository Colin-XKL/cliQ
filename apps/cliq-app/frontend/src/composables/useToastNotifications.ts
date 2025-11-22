import { useToast } from "primevue/usetoast";
import Toast from "primevue/toast";
import ToastService from "primevue/toastservice";
type ToastLevel = "success" | "info" | "warn" | "error";

export function useToastNotifications() {
  const toast = useToast();

  const showToast = (summary: string, detail: string, severity: ToastLevel) => {
    toast.add({
      severity,
      summary,
      detail,
      life: 3000,
    });
  };

  return { showToast };
}
