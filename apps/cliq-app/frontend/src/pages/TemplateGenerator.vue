<template>
  <div class="p-6 max-w-6xl mx-auto">
    <div class="mb-6">
      <div class="flex items-center gap-4 mb-2">
        <span class="text-sm text-gray-600">模式</span>
        <Select v-model="activeMode" :options="modeOptions" optionLabel="label" optionValue="value" class="w-64" />
      </div>
      <div class="text-xs text-gray-500">当前模式：{{ activeModeLabel }}</div>
    </div>

    <!-- 输入区域 -->
    <Transition name="fade" mode="out-in">
      <div v-if="activeMode === 'custom'" key="custom" class="mb-6">
        <InputText v-model="commandInput" type="text"
          placeholder="输入命令例如: pngquant {{input_file}} --output {{output_file}}"
          class="w-full p-3 border border-gray-300 rounded-md" />
        <div class="flex gap-3 my-6">
          <Button @click="generateTemplate" :label="isGenerating ? '生成中...' : '开始转换'" :disabled="isGenerating"
            class="bg-purple-500 hover:bg-purple-600 text-white" />
        </div>
      </div>
      <div v-else key="smart" class="mb-6">
        <div class="space-y-4">
          <InputText v-model="smartCommandExample" type="text" placeholder="示例命令，如: pngquant input.png --output output.png"
            class="w-full p-3 border border-gray-300 rounded-md" />
          <InputText v-model="smartDescription" type="text" placeholder="自然语言说明（可选）"
            class="w-full p-3 border border-gray-300 rounded-md" />
        </div>
        <div class="flex gap-3 my-6">
          <Button @click="smartGenerateTemplate" :label="isGenerating ? '生成中...' : '生成模板'" :disabled="isGenerating"
            class="bg-purple-500 hover:bg-purple-600 text-white" />
        </div>
        <div v-if="inputError" class="text-sm text-red-600">{{ inputError }}</div>
      </div>
    </Transition>

    <!-- 模板输出区域 -->
    <div>
      <Card>
        <template #title>
          <h2>生成的模板</h2>
          <div class="flex items-center justify-between">
            <div class="flex gap-2">
              <Button @click="copyToClipboard" severity="secondary" label="复制" icon="pi pi-copy" size="small" />
              <Button @click="exportTemplate" severity="secondary" label="导出" icon="pi pi-download" size="small" />
              <Button @click="saveToLocal" severity="secondary" label="收藏保存" icon="pi pi-bookmark" size="small" />
              <Button @click="openTemplateEditor" severity="secondary" label="编辑模板" icon="pi pi-pencil" size="small" />
            </div>
          </div>
        </template>
        <template #content>
          <div class="w-full">
            <textarea v-model="generatedYaml" readonly
              class="w-full h-96 p-3 border border-gray-300 rounded-md font-mono text-sm resize-none"
              placeholder="生成的模板将显示在这里..."></textarea>
            <div v-if="validationError" class="mt-2 text-sm text-red-600">{{ validationError }}</div>
          </div>
        </template>
      </Card>
    </div>

    <Transition name="fade" mode="out-in">
      <div v-if="parsedTemplate" key="preview" class="mt-6 grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="space-y-4">
          <TemplateMetadataDisplay :template="parsedTemplate" />
          <Card>
            <template #title>
              <div class="flex items-center justify-between">
                <h3>参数预览</h3>
              </div>
            </template>
            <template #content>
              <DynamicCommandForm :selectedCommand="selectedPreviewCommand" v-model:commandVariableValues="previewVars" />
            </template>
          </Card>
        </div>
        <div class="space-y-4">
          <Card>
            <template #title>
              <div class="flex items-center justify-between">
                <h3>命令预览</h3>
              </div>
            </template>
            <template #content>
              <pre class="p-3 bg-gray-50 rounded-md text-sm break-words">{{ previewCommandText }}</pre>
            </template>
          </Card>
        </div>
      </div>
    </Transition>

    <!-- Template Editor Modal -->
    <TemplateEditorModal :visible="showEditorModal" :initialYaml="generatedYaml" @close="showEditorModal = false"
      @save="onTemplateSaved" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { ParseCommandToTemplate, GenerateYAMLFromTemplate, SaveYAMLToFile, ParseYAMLToTemplate, SaveFavTemplate, ValidateYAMLTemplate } from '@/wailsjs/go/main/App';
import { useToastNotifications } from '@/composables/useToastNotifications';
import TemplateEditorModal from '@/components/TemplateEditorModal.vue';
import { useSettings, DEFAULT_BASE_URL } from '@/composables/useSettings';

const { showToast } = useToastNotifications();

const commandInput = ref('');
const generatedYaml = ref('');
const showEditorModal = ref(false);

const activeMode = ref<'custom' | 'smart'>('custom');
const modeOptions = [{ label: '自定义模式', value: 'custom' }, { label: '智能生成模式', value: 'smart' }];
const activeModeLabel = computed(() => modeOptions.find(o => o.value === activeMode.value)?.label || '');

const smartCommandExample = ref('');
const smartDescription = ref('');
const isGenerating = ref(false);
const inputError = ref('');
const validationError = ref('');
const parsedTemplate = ref<any | null>(null);
const selectedPreviewCommand = ref<any | null>(null);
const previewVars = ref<Record<string, any>>({});
const previewCommandText = computed(() => {
  if (!selectedPreviewCommand.value || !selectedPreviewCommand.value.command) return '';
  let cmd = selectedPreviewCommand.value.command as string;
  return cmd.replace(/\{\{\s*([\w-]+)\s*\}\}/g, (_, name) => {
    const v = previewVars.value[name];
    return v !== undefined && v !== null ? String(v) : '';
  });
});

const generateTemplate = async () => {
  if (!commandInput.value.trim()) {
    showToast('错误', '请输入CLI命令', 'error');
    return;
  }

  try {
    const templateObj = await ParseCommandToTemplate(commandInput.value);
    if (templateObj) {
      const yamlStr = await GenerateYAMLFromTemplate(templateObj);
      generatedYaml.value = yamlStr;

      showToast('成功', '模板生成成功', 'success');
    }
  } catch (error) {
    showToast('错误', `生成模板失败: ${error}`, 'error');
    console.error('生成模板失败:', error);
  }
};

const smartGenerateTemplate = async () => {
  inputError.value = '';
  if (!smartCommandExample.value.trim()) {
    inputError.value = '示例命令为必填项';
    showToast('错误', inputError.value, 'error');
    return;
  }
  try {
    isGenerating.value = true;
    const { loadSettings } = useSettings();
    const s = await loadSettings();
    const base = (s.cliq_hub_base_url || DEFAULT_BASE_URL).replace(/\/$/, '');
    const url = `${base}/v1/templates/generate`;
    const resp = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command_example: smartCommandExample.value, description: smartDescription.value, encoding: 'plain' })
    });
    if (!resp.ok) {
      const text = await resp.text();
      throw new Error(text || `HTTP ${resp.status}`);
    }
    const data = await resp.json();
    let yamlStr: string = data.yaml || '';
    if (data.encoding === 'base64' && yamlStr) {
      yamlStr = atob(yamlStr);
    }
    if (!yamlStr) {
      throw new Error('后端未返回有效的模板内容');
    }
    generatedYaml.value = yamlStr;
    showToast('成功', '模板生成成功', 'success');
  } catch (e: any) {
    showToast('错误', `智能生成失败: ${e}`, 'error');
    console.error('智能生成失败:', e);
  } finally {
    isGenerating.value = false;
  }
};

const copyToClipboard = async () => {
  if (!generatedYaml.value) {
    showToast('错误', '没有可复制的模板', 'error');
    return;
  }

  try {
    await navigator.clipboard.writeText(generatedYaml.value);
    showToast('成功', '模板已复制到剪贴板', 'success');
  } catch (error) {
    showToast('错误', '复制失败: ' + error, 'error');
  }
};

const exportTemplate = async () => {
  if (!generatedYaml.value) {
    showToast('错误', '没有可导出的模板', 'error');
    return;
  }

  try {
    await SaveYAMLToFile(generatedYaml.value);
    showToast('成功', '模板已导出', 'success');
  } catch (error) {
    showToast('错误', '导出失败: ' + error, 'error');
  }
};

const saveToLocal = async () => {
  if (!generatedYaml.value) {
    showToast('错误', '没有可保存的模板', 'error');
    return;
  }

  try {
    // Parse the YAML back to a template object
    const templateObj = await ParseYAMLToTemplate(generatedYaml.value);

    // Ensure the template has a name for saving
    if (!templateObj.name) {
      // Generate a name based on the command if not provided
      const commandName = commandInput.value.trim().split(' ')[0] || 'unnamed-template';
      templateObj.name = `${commandName}-${Date.now()}`;
    }

    // Save the template to local storage via backend
    await SaveFavTemplate(templateObj);
    showToast('成功', `模板 "${templateObj.name}" 已收藏保存`, 'success');
  } catch (error) {
    showToast('错误', '保存模板失败: ' + error, 'error');
    console.error('保存模板失败:', error);
  }
};

const openTemplateEditor = () => {
  if (!generatedYaml.value) {
    showToast('提示', '请先生成模板', 'warn');
    return;
  }
  showEditorModal.value = true;
};

const onTemplateSaved = (updatedYaml: string) => {
  generatedYaml.value = updatedYaml;
  showToast('成功', '模板已更新', 'success');
};

watch(generatedYaml, async (val) => {
  validationError.value = '';
  parsedTemplate.value = null;
  selectedPreviewCommand.value = null;
  previewVars.value = {};
  if (!val) return;
  try {
    await ValidateYAMLTemplate(val);
    const tpl = await ParseYAMLToTemplate(val);
    parsedTemplate.value = tpl;
    if (tpl && tpl.cmds && tpl.cmds.length > 0) {
      selectedPreviewCommand.value = tpl.cmds[0];
      const vars: Record<string, any> = {};
      (selectedPreviewCommand.value.variables || []).forEach((v: any) => { vars[v.name] = ''; });
      previewVars.value = vars;
    }
  } catch (err: any) {
    validationError.value = String(err);
  }
});
</script>