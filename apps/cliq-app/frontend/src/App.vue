<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { models } from '@/wailsjs/go/models';
import MainPage from '@/pages/MainPage.vue';
import { DynamicCommandForm } from '@repo/shared-vue-ui';
import CommandExecutor from '@/components/CommandExecutor.vue';
import TemplateGenerator from '@/pages/TemplateGenerator.vue';
import TemplateManagementPage from '@/pages/TemplateManagementPage.vue';
import AboutPage from '@/pages/AboutPage.vue';
import SettingsPage from '@/pages/SettingsPage.vue';

declare global {
  interface Window {
    go: any;
  }
}

const templateData = ref<models.TemplateFile>({} as models.TemplateFile);
const selectedCommand = ref<any>(null);
const commandVariableValues = ref<{ [key: string]: any }>({});
const isProcessing = ref(false);
const commandOutput = ref('');
const currentView = ref<'main' | 'generator' | 'template-management' | 'about' | 'settings'>('main'); // Add view state
const favTemplates = ref<models.TemplateFile[]>([]);

const resetTemplate = () => {
  templateData.value = {} as models.TemplateFile;
  selectedCommand.value = null;
  commandVariableValues.value = {};
  isProcessing.value = false;
  commandOutput.value = '';
};

const loadFavTemplates = async () => {
  try {
    if (window.go && window.go.main && window.go.main.App && window.go.main.App.ListFavTemplates) {
      const result = await window.go.main.App.ListFavTemplates();
      favTemplates.value = result || [];
    } else {
      console.warn('Wails backend not available. Skipping loading favorite templates.');
      favTemplates.value = [];
    }
  } catch (error) {
    console.error('Failed to list favorite templates:', error);
    favTemplates.value = [];
  }
};

onMounted(() => {
  loadFavTemplates();
});
</script>

<template>
  <div class="homepage-bg h-full w-full">
    <div class="flex flex-col items-center justify-center h-[100]vh p-4">
      <div class="w-full max-w-5xl">
        <!-- Navigation tabs -->
        <div class="flex justify-center mb-6">
          <div class="inline-flex bg-gray-100 rounded-lg p-1">
            <button @click="currentView = 'main'"
              :class="['px-4 py-2 rounded-md text-sm font-medium', currentView === 'main' ? 'bg-white shadow text-gray-900' : 'text-gray-600']">
              主界面
            </button>
            <button @click="currentView = 'generator'"
              :class="['px-4 py-2 rounded-md text-sm font-medium', currentView === 'generator' ? 'bg-white shadow text-gray-900' : 'text-gray-600']">
              模板生成器
            </button>
            <button @click="currentView = 'template-management'"
              :class="['px-4 py-2 rounded-md text-sm font-medium', currentView === 'template-management' ? 'bg-white shadow text-gray-900' : 'text-gray-600']">
              模板管理
            </button>
            <button @click="currentView = 'about'"
              :class="['px-4 py-2 rounded-md text-sm font-medium', currentView === 'about' ? 'bg-white shadow text-gray-900' : 'text-gray-600']">
              关于
            </button>
            <button @click="currentView = 'settings'"
              :class="['px-4 py-2 rounded-md text-sm font-medium', currentView === 'settings' ? 'bg-white shadow text-gray-900' : 'text-gray-600']">
              设置
            </button>
          </div>
        </div>

        <div class="bg-white p-6 rounded-lg shadow-md overflow-y-auto max-h-[80vh]">
          <!-- Main View -->
          <div v-if="currentView === 'main'">
            <MainPage v-model:templateData="templateData" v-model:selectedCommand="selectedCommand"
              @reset-template="resetTemplate" :favTemplates="favTemplates" @fav-template-updated="loadFavTemplates" />

            <DynamicCommandForm v-if="templateData.name" :selectedCommand="selectedCommand"
              v-model:commandVariableValues="commandVariableValues" />

            <CommandExecutor v-if="templateData.name" :selectedCommand="selectedCommand"
              :commandVariableValues="commandVariableValues" v-model:isProcessing="isProcessing"
              v-model:commandOutput="commandOutput" />
          </div>

          <!-- Template Generator View -->
          <div v-if="currentView === 'generator'">
            <TemplateGenerator />
          </div>

          <!-- Template Management View -->
          <div v-if="currentView === 'template-management'">
            <TemplateManagementPage />
          </div>

          <!-- About View -->
          <div v-if="currentView === 'about'">
            <AboutPage />
          </div>

          <!-- Settings View -->
          <div v-if="currentView === 'settings'">
            <SettingsPage />
          </div>
        </div>
      </div>
    </div>
  </div>
  <Toast />
</template>

<style>
@import './styles/home-bg.css';
</style>
