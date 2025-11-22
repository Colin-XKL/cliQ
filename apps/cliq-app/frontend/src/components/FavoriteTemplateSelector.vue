<template>
  <div>
    <button @click="openDialog"
      class="bg-indigo-500 text-white px-4 py-2 rounded-md hover:bg-indigo-600 focus:outline-none">
      从收藏夹选择
    </button>

    <!-- Favorite Template Selection Dialog -->
    <div v-if="showDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded-lg w-full max-w-md max-h-90vh overflow-y-auto">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-semibold">从收藏夹选择模板</h3>
          <button @click="closeDialog" class="text-gray-500 hover:text-gray-700">
            <i class="pi pi-times"></i>
          </button>
        </div>

        <div v-if="favTemplates && favTemplates.length > 0" class="space-y-2">
          <div v-for="template in favTemplates" :key="template.name"
            class="p-4 border rounded-lg shadow-sm cursor-pointer hover:bg-gray-100 border-gray-300"
            @click="selectTemplate(template.name)">
            <h4 class="font-semibold text-gray-800">{{ template.name }}</h4>
            <p class="text-sm text-gray-500 truncate">{{ template.description }}</p>
          </div>
        </div>

        <div v-else class="text-center py-8 text-gray-500">
          <p>暂无收藏模板</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { models } from '@/wailsjs/go/models';
import { GetFavTemplate } from '@/wailsjs/go/main/App';
import { useToastNotifications } from '@/composables/useToastNotifications';

interface Props {
  favTemplates: models.TemplateFile[];
}

interface Emits {
  (e: 'template-selected', template: models.TemplateFile): void;
  (e: 'close'): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const { showToast } = useToastNotifications();
const showDialog = ref(false);

const openDialog = () => {
  showDialog.value = true;
};

const closeDialog = () => {
  showDialog.value = false;
  emit('close');
};

const selectTemplate = async (templateName: string) => {
  try {
    const result = await GetFavTemplate(templateName);
    if (result) {
      emit('template-selected', result);
      closeDialog();
      showToast('成功', `模板 ${templateName} 加载成功`, 'success');
    }
  } catch (error) {
    showToast('错误', `加载收藏模板失败: ${error}`, 'error');
    console.error('加载收藏模板失败:', error);
  }
};
</script>