import { ref } from 'vue'

export type AppSettings = {
  cliq_hub_base_url: string
}

export const DEFAULT_BASE_URL = 'http://localhost:8080'

const settingsRef = ref<AppSettings | null>(null)

export const useSettings = () => {
  const loadSettings = async (): Promise<AppSettings> => {
    try {
      if (window.go && window.go.main && window.go.main.App && window.go.main.App.GetAppSettings) {
        const s = await window.go.main.App.GetAppSettings()
        settingsRef.value = s || { cliq_hub_base_url: DEFAULT_BASE_URL }
      } else {
        settingsRef.value = { cliq_hub_base_url: DEFAULT_BASE_URL }
      }
    } catch {
      settingsRef.value = { cliq_hub_base_url: DEFAULT_BASE_URL }
    }
    return settingsRef.value as AppSettings
  }

  const saveSettings = async (partial: Record<string, any>): Promise<void> => {
    if (window.go && window.go.main && window.go.main.App && window.go.main.App.UpdateAppSettings) {
      await window.go.main.App.UpdateAppSettings(partial)
      settingsRef.value = { ...(settingsRef.value || { cliq_hub_base_url: DEFAULT_BASE_URL }), ...partial } as AppSettings
    } else {
      // No-op if backend bindings not available
      settingsRef.value = { ...(settingsRef.value || { cliq_hub_base_url: DEFAULT_BASE_URL }), ...partial } as AppSettings
    }
  }

  return { settings: settingsRef, loadSettings, saveSettings }
}