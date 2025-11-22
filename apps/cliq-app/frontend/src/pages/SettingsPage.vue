<template>
  <div class="p-6 max-w-3xl mx-auto">
    <h2 class="text-xl font-semibold mb-4">设置</h2>
    <Card>
      <template #title>
        <div class="flex items-center justify-between">
          <span>CliQ Hub Base URL</span>
        </div>
      </template>
      <template #content>
        <div class="space-y-3">
          <InputText v-model="baseUrl" type="text" :placeholder="DEFAULT_BASE_URL"
            class="w-full p-3 border border-gray-300 rounded-md" />
          <div v-if="error" class="text-sm text-red-600">{{ error }}</div>
          <div class="flex gap-3 mt-4">
            <Button :disabled="saving || !!error" @click="onSave" class="bg-purple-500 hover:bg-purple-600 text-white"
              :label="saving ? '保存中...' : '保存'" />
            <Button :disabled="saving" @click="onReset" severity="secondary" label="重置为默认" />
          </div>
        </div>
      </template>
    </Card>
  </div>
  
  <Toast />
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useSettings, DEFAULT_BASE_URL } from '@/composables/useSettings'
import { useToastNotifications } from '@/composables/useToastNotifications'

const { showToast } = useToastNotifications()
const { settings, loadSettings, saveSettings } = useSettings()

const baseUrl = ref('')
const error = ref('')
const saving = ref(false)

const validate = (val: string) => {
  error.value = ''
  try {
    const parsed = new URL(val)
    if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') {
      throw new Error('仅支持 http 或 https 协议')
    }
  } catch (e: any) {
    error.value = 'URL 格式不正确'
  }
}

watch(baseUrl, (val) => validate(val))

onMounted(async () => {
  const s = await loadSettings()
  baseUrl.value = s.cliq_hub_base_url || DEFAULT_BASE_URL
})

const onSave = async () => {
  if (error.value) return
  try {
    saving.value = true
    await saveSettings({ cliq_hub_base_url: baseUrl.value.replace(/\/$/, '') })
    showToast('成功', '配置已保存', 'success')
  } catch (e: any) {
    showToast('错误', String(e), 'error')
  } finally {
    saving.value = false
  }
}

const onReset = async () => {
  try {
    saving.value = true
    baseUrl.value = DEFAULT_BASE_URL
    await saveSettings({ cliq_hub_base_url: DEFAULT_BASE_URL })
    showToast('成功', '已重置为默认', 'success')
  } catch (e: any) {
    showToast('错误', String(e), 'error')
  } finally {
    saving.value = false
  }
}
</script>