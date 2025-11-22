<template>
  <div v-if="!templateDataInternal.name" class="text-center py-8">
    <div class="text-center mb-8">
      <h2 class="text-2xl font-bold text-black mb-4">欢迎使用 cliQ</h2>
      <p class="text-xl mt-2">将复杂的 CLI 命令转化为直观、易用的图形用户界面</p>
    </div>
    <p class="text-gray-400 mb-4">请导入模板文件以开始使用</p>
    <div class="flex justify-center gap-4">
      <button @click="importTemplate"
        class="bg-purple-500 text-white px-6 py-3 rounded-md hover:bg-purple-600 focus:outline-none text-lg">
        从文件导入
      </button>
      <button @click="showUrlImportDialog = true"
        class="bg-blue-500 text-white px-6 py-3 rounded-md hover:bg-blue-600 focus:outline-none text-lg">
        从URL导入
      </button>
    </div>
    <div v-if="favTemplates && favTemplates.length > 0" class="mt-8">
      <p class="text-gray-400 mb-4">或从收藏夹选择</p>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="(template, index) in favTemplates.slice(0, 9)" :key="template.name"
          class="p-4 border rounded-lg shadow-sm cursor-pointer hover:bg-gray-100 border-gray-300"
          @click="loadFavoriteTemplate(template.name)">
          <h4 class="font-semibold text-gray-800">{{ template.name }}</h4>
          <p class="text-sm text-gray-500 truncate">{{ template.description }}</p>
        </div>
      </div>
    </div>

  </div>

  <div v-else>
    <div class="flex justify-between items-center mb-4">
      <div class="flex gap-2">
        <button @click="importTemplate"
          class="bg-purple-500 text-white px-4 py-2 rounded-md hover:bg-purple-600 focus:outline-none">
          从文件导入
        </button>
        <FavoriteTemplateSelector :favTemplates="favTemplates" @template-selected="handleTemplateSelected"
          @close="handleSelectorClose" />
        <button @click="showUrlImportDialog = true"
          class="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 focus:outline-none">
          URL导入
        </button>
        <button @click="addTemplateToFavorites"
          class="bg-yellow-500 text-white px-4 py-2 rounded-md hover:bg-yellow-600 focus:outline-none">
          收藏
        </button>
      </div>
    </div>

    <!-- 模板基本信息显示 -->
    <TemplateMetadataDisplay :template="templateDataInternal" class="mb-6" />

    <!-- 命令选择 -->
    <div class="mb-6" v-if="templateDataInternal.cmds && templateDataInternal.cmds.length > 0">
      <label class="block text-sm font-medium text-gray-700 mb-2">选择命令</label>
      <Dropdown v-model="selectedCommandInternal" :options="templateDataInternal.cmds" optionLabel="name" class="w-full"
        placeholder="选择要执行的命令">
        <template #value="slotProps">
          <div class="flex align-items-center">
            <div>{{ slotProps.value.name }}</div>
          </div>
        </template>
        <template #option="slotProps">
          <div class="flex flex-col text-left">
            <div>{{ slotProps.option.name }}</div>
          </div>
        </template>
      </Dropdown>
      <p class="mt-2 text-sm text-gray-500" v-if="selectedCommandInternal">{{ selectedCommandInternal.description }}</p>
    </div>
  </div>

  <!-- URL导入对话框 -->
  <div v-if="showUrlImportDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white p-6 rounded-lg w-full max-w-md">
      <h3 class="text-lg font-semibold mb-4">从URL导入模板</h3>
      <div class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">模板URL</label>
        <input v-model="templateUrl" type="text"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="https://example.com/template.cliqfile.yaml">
      </div>
      <div class="flex justify-end gap-2">
        <button @click="cancelUrlImport" class="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md">
          取消
        </button>
        <button @click="importTemplateFromUrl" class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600">
          导入
        </button>
      </div>
    </div>
  </div>

</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import { ImportTemplate, ImportTemplateFromURL, GetFavTemplate } from '@/wailsjs/go/main/App';
import { models } from '@/wailsjs/go/models';
import { useToastNotifications } from '@/composables/useToastNotifications';
import { SaveFavTemplate } from '@/wailsjs/go/main/App';
import TemplateMetadataDisplay from '@/components/TemplateMetadataDisplay.vue';
import FavoriteTemplateSelector from '@/components/FavoriteTemplateSelector.vue';

const props = defineProps({
  templateData: { type: Object as () => models.TemplateFile, required: true },
  selectedCommand: { type: Object as () => any, default: null },
  favTemplates: { type: Array as () => models.TemplateFile[], default: () => [] },
});

const emit = defineEmits(['update:templateData', 'update:selectedCommand', 'reset-template', 'fav-template-updated']);

const { showToast } = useToastNotifications();

const templateDataInternal = ref(props.templateData);
const selectedCommandInternal = ref(props.selectedCommand);
const showUrlImportDialog = ref(false);
const templateUrl = ref('');

// Helper function to update template state consistently
const updateTemplateState = (template: models.TemplateFile) => {
  // Emit reset to clear any existing state in related components
  emit('reset-template');

  // Update internal state
  templateDataInternal.value = template;
  selectedCommandInternal.value = null; // Reset selected command on new template import
  if (template.cmds && template.cmds.length > 0) {
    selectedCommandInternal.value = template.cmds[0];
  }

  // Important: After parent state is reset, update it with the new template data
  // This ensures that the parent's reactive system picks up the changes
  emit('update:templateData', template);
  emit('update:selectedCommand', selectedCommandInternal.value);

  // Force next tick update to ensure all components re-render properly
  // This ensures that child components like DynamicCommandForm get updated
  setTimeout(() => {
    // No additional action needed, the emits should handle the updates
  }, 0);
};

watch(() => props.templateData, (newValue) => {
  templateDataInternal.value = newValue;
});

watch(() => props.selectedCommand, (newValue) => {
  selectedCommandInternal.value = newValue;
});

watch(templateDataInternal, (newValue) => {
  emit('update:templateData', newValue);
}, { deep: true });

watch(selectedCommandInternal, (newValue) => {
  emit('update:selectedCommand', newValue);
});

const importTemplate = async () => {
  try {
    emit('reset-template');
    const result = await ImportTemplate();
    if (result) {
      templateDataInternal.value = result;
      selectedCommandInternal.value = null; // Reset selected command on new template import
      if (result.cmds && result.cmds.length > 0) {
        selectedCommandInternal.value = result.cmds[0];
      }
      showToast('成功', '模板导入成功', 'success');
    }
  } catch (error) {
    showToast('错误', `导入模板失败: ${error}`, 'error');
    console.error('导入模板失败:', error);
  }
};

const importTemplateFromUrl = async () => {
  if (!templateUrl.value.trim()) {
    showToast('错误', '请输入模板URL', 'error');
    return;
  }

  try {
    emit('reset-template');
    const result = await ImportTemplateFromURL(templateUrl.value);
    if (result) {
      templateDataInternal.value = result;
      selectedCommandInternal.value = null; // Reset selected command on new template import
      if (result.cmds && result.cmds.length > 0) {
        selectedCommandInternal.value = result.cmds[0];
      }
      showToast('成功', '模板导入成功', 'success');
      cancelUrlImport(); // Close the dialog
    }
  } catch (error) {
    showToast('错误', `从URL导入模板失败: ${error}`, 'error');
    console.error('从URL导入模板失败:', error);
  }
};



const cancelUrlImport = () => {
  showUrlImportDialog.value = false;
  templateUrl.value = '';
};

const addTemplateToFavorites = async () => {
  if (!templateDataInternal.value || !templateDataInternal.value.name) {
    showToast('错误', '没有可收藏的模板', 'error');
    return;
  }

  try {
    await SaveFavTemplate(templateDataInternal.value);
    showToast('成功', `模板 ${templateDataInternal.value.name} 已收藏`, 'success');
    emit('fav-template-updated'); // Notify parent to refresh favorite templates
  } catch (error) {
    showToast('错误', `收藏模板失败: ${error}`, 'error');
    console.error('收藏模板失败:', error);
  }
};

const loadFavoriteTemplate = async (templateName: string) => {
  try {
    const result = await GetFavTemplate(templateName);
    if (result) {
      updateTemplateState(result);
      showToast('成功', `模板 ${templateName} 加载成功`, 'success');
    }
  } catch (error) {
    showToast('错误', `加载收藏模板失败: ${error}`, 'error');
    console.error('加载收藏模板失败:', error);
  }
};

const handleTemplateSelected = (template: models.TemplateFile) => {
  try {
    updateTemplateState(template);
  } catch (error) {
    showToast('错误', `处理收藏模板失败: ${error}`, 'error');
    console.error('处理收藏模板失败:', error);
  }
};

const handleSelectorClose = () => {
  // No specific action needed when selector closes
};

</script>