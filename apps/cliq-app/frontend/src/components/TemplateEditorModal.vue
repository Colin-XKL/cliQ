<template>
  <div v-if="visible" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white p-6 rounded-lg w-full max-w-6xl h-5/6 flex flex-col">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-xl font-semibold">模板编辑器</h3>
        <Button @click="closeModal" label="关闭" class="p-button-secondary" />
      </div>

      <div class="flex flex-col flex-grow min-h-0">
        <div class="flex gap-4 mb-4">
          <Button @click="validateTemplate" label="校验模板" />
          <Button @click="applyChanges" label="应用更改" class="p-button-success" />
        </div>

        <div class="flex flex-col lg:flex-row gap-6 flex-grow min-h-0">
          <!-- 编辑器区域 -->
          <div class="flex flex-col flex-1 min-h-0">
            <h4 class="text-lg font-medium mb-2">模板编辑</h4>
            <div class="flex-grow border border-gray-300 rounded-md overflow-y-auto">
              <MonacoEditor :value="templateYaml" :key="editorKey" language="yaml" theme="vs" :options="{
                minimap: { enabled: false },
                automaticLayout: true,
                fontSize: 14,
                scrollBeyondLastLine: false,
                wordWrap: 'on',
                formatOnType: true,
                formatOnPaste: true,
                readOnly: false
              }" @change="onEditorChange" />
            </div>
          </div>

          <!-- 预览区域 -->
          <div class="flex flex-col flex-1 min-h-0">
            <h4 class="text-lg font-medium mb-2">表单预览</h4>
            <div class="flex-grow p-4 bg-gray-50 rounded-md overflow-y-auto">
              <!-- Template metadata display -->
              <div class="mb-4" v-if="fullTemplateData">
                <TemplateMetadataDisplay :template="fullTemplateData" />
              </div>

              <!-- 命令选择下拉框 -->
              <div class="mb-4" v-if="fullTemplateData && fullTemplateData.cmds && fullTemplateData.cmds.length > 0">
                <label class="block text-sm font-medium text-gray-700 mb-2">选择命令</label>
                <Dropdown v-model="selectedPreviewCommand" :options="fullTemplateData.cmds" optionLabel="name"
                  class="w-full" placeholder="选择要预览的命令">
                  <template #value="slotProps">
                    <div class="flex align-items-center">
                      <div>{{ slotProps.value?.name || '选择命令' }}</div>
                    </div>
                  </template>
                  <template #option="slotProps">
                    <div class="flex flex-col text-left">
                      <div>{{ slotProps.option.name }}</div>
                      <small v-if="slotProps.option.description" class="text-gray-500">{{ slotProps.option.description
                      }}</small>
                    </div>
                  </template>
                </Dropdown>
              </div>

              <div v-if="selectedPreviewCommand" class="w-full">
                <DynamicCommandForm :selectedCommand="selectedPreviewCommand"
                  :commandVariableValues="commandVariableValues"
                  @update:commandVariableValues="updateCommandVariableValues" />
              </div>
              <div v-else-if="hasValidationError" class="flex flex-col items-center justify-center h-64 text-red-500">
                <i class="pi pi-exclamation-triangle text-4xl mb-3"></i>
                <p>模板格式无效，请检查YAML语法</p>
              </div>
              <div v-else class="flex items-center justify-center h-64 text-gray-500">
                <p>校验模板后将显示表单预览</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import MonacoEditor from 'monaco-editor-vue3';
import { DynamicCommandForm } from '@repo/shared-vue-ui';
import TemplateMetadataDisplay from '@/components/TemplateMetadataDisplay.vue';
import { ValidateYAMLTemplate, ParseYAMLToTemplate } from '@/wailsjs/go/main/App';
import { models } from '@/wailsjs/go/models';
import { useToastNotifications } from '@/composables/useToastNotifications';
import Dropdown from 'primevue/dropdown';

const { showToast } = useToastNotifications();

interface Props {
  visible: boolean;
  initialYaml: string;
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  initialYaml: ''
});

const emit = defineEmits(['close', 'save']);

const templateYaml = ref(props.initialYaml);
const editorKey = ref(0);
const previewCommand = ref<any>(null);
const fullTemplateData = ref<models.TemplateFile | null>(null);
const selectedPreviewCommand = ref<any>(null);
const hasValidationError = ref(false);
const commandVariableValues = reactive<{ [key: string]: any }>({});

// Track if the change is coming from the editor to avoid infinite loops
let isUpdatingFromEditor = false;

const onEditorChange = (value: string) => {
  isUpdatingFromEditor = true;
  templateYaml.value = value;
  // The watch handler will take care of the preview update
};

// Watch for changes in templateYaml and update preview accordingly
watch(templateYaml, async (newYaml) => {
  if (!isUpdatingFromEditor) {
    // This means the value was updated from outside the editor
    // Force refresh the editor by updating the key
    editorKey.value += 1;
  } else {
    // Reset the flag after delay to ensure the editor state is stable
    setTimeout(() => {
      isUpdatingFromEditor = false;
    }, 0);
  }

  if (newYaml) {
    await updatePreviewFromYaml();
  } else {
    previewCommand.value = null;
    fullTemplateData.value = null;
    selectedPreviewCommand.value = null;
    hasValidationError.value = false;
  }
});

watch(() => props.visible, async (newVisible) => {
  if (newVisible && props.initialYaml) {
    // Always reset to initialYaml when modal becomes visible
    templateYaml.value = props.initialYaml;
    // Force refresh the editor by updating the key to ensure Monaco picks up the value
    editorKey.value += 1;
    // Give Monaco a moment to reinitialize, then update the preview
    setTimeout(async () => {
      await updatePreviewFromYaml();
    }, 100);
  }
});

// Watch for changes in selectedPreviewCommand to update commandVariableValues
watch(selectedPreviewCommand, (newCommand) => {
  if (newCommand && newCommand.variables) {
    // Clear the current commandVariableValues
    Object.keys(commandVariableValues).forEach(key => {
      delete commandVariableValues[key];
    });

    // Initialize command variable values for the selected command
    // Check if variables is an array of VariableDefinition objects (with flattened structure)
    if (Array.isArray(newCommand.variables)) {
      (newCommand.variables as models.VariableDefinition[]).forEach(varDef => {
        commandVariableValues[varDef.name] = undefined; // Initialize with undefined
      });
    } else {
      // Fallback for old map format
      Object.keys(newCommand.variables).forEach(key => {
        commandVariableValues[key] = undefined; // Initialize with undefined
      });
    }
  } else {
    // If no command is selected, clear all variables
    Object.keys(commandVariableValues).forEach(key => {
      delete commandVariableValues[key];
    });
  }
});

const updatePreview = async (templateObj: models.TemplateFile) => {
  // Store the full template data for metadata display
  fullTemplateData.value = templateObj;

  if (templateObj && templateObj.cmds && templateObj.cmds.length > 0) {
    // Use the first command for preview by default
    previewCommand.value = templateObj.cmds[0];
    selectedPreviewCommand.value = templateObj.cmds[0];

    // Initialize command variable values
    if (previewCommand.value.variables) {
      // Check if variables is an array of VariableDefinition objects (with flattened structure)
      if (Array.isArray(previewCommand.value.variables)) {
        (previewCommand.value.variables as models.VariableDefinition[]).forEach(varDef => {
          commandVariableValues[varDef.name] = undefined; // Initialize with undefined
        });
      } else {
        // Fallback for old map format
        Object.keys(previewCommand.value.variables).forEach(key => {
          commandVariableValues[key] = undefined; // Initialize with undefined
        });
      }
    }
  } else {
    // Reset if there are no commands
    previewCommand.value = null;
    selectedPreviewCommand.value = null;
  }
};

const updatePreviewFromYaml = async () => {
  if (!templateYaml.value) {
    previewCommand.value = null;
    fullTemplateData.value = null;
    hasValidationError.value = false;
    return;
  }

  try {
    // Parse the YAML back to a template object for preview
    const templateObj = await ParseYAMLToTemplate(templateYaml.value);
    updatePreview(templateObj);
    hasValidationError.value = false; // Clear any validation errors
  } catch (error) {
    console.error('解析YAML模板失败:', error);
    previewCommand.value = null;
    fullTemplateData.value = null;
    hasValidationError.value = true; // Set validation error state
  }
};

const updateCommandVariableValues = (newValues: { [key: string]: any }) => {
  Object.assign(commandVariableValues, newValues);
};

const validateTemplate = async () => {
  if (!templateYaml.value) {
    showToast('错误', '没有可校验的模板', 'error');
    return;
  }

  try {
    await ValidateYAMLTemplate(templateYaml.value);
    showToast('成功', '模板格式有效', 'success');
    hasValidationError.value = false; // Clear validation error state

    // Update preview after validation since it's a valid template
    await updatePreviewFromYaml();
  } catch (error) {
    showToast('错误', '模板格式无效: ' + error, 'error');
    hasValidationError.value = true; // Set validation error state
    selectedPreviewCommand.value = null; // Reset selected command on validation error
  }
};




const applyChanges = () => {
  emit('save', templateYaml.value);
  closeModal();
};

const closeModal = () => {
  emit('close');
};

</script>