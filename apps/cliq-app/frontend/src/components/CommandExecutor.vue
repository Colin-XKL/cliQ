<template>
  <!-- 运行按钮 -->
  <div class="flex justify-center space-x-4">
    <button @click="runCommand"
      class="bg-green-500 text-white px-6 py-3 rounded-md hover:bg-green-600 focus:outline-none disabled:bg-gray-400"
      :disabled="isProcessingInternal">
      <span v-if="isProcessingInternal">处理中...</span>
      <span v-else>运行命令</span>
    </button>
    <button @click="showCommandInfo"
      class="bg-blue-500 text-white px-6 py-3 rounded-md hover:bg-blue-600 focus:outline-none">
      查看命令
    </button>
  </div>

  <!-- 执行结果模态框 -->
  <div v-if="showResultModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
    @click.self="closeResultModal">
    <div
      class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white max-w-4xl max-h-[80vh] overflow-y-auto">
      <div class="mt-3 text-center">
        <!-- 执行状态指示器 -->
        <div v-if="executionState === 'executing'">
          <div class="flex items-center justify-center">
            <i class="pi pi-spin pi-spinner text-4xl text-blue-500 mr-4"></i>
            <h3 class="text-lg leading-6 font-medium text-blue-600">
              命令执行中...
            </h3>
          </div>
        </div>

        <!-- 完成状态指示器 -->
        <div v-else-if="executionState === 'completed'"
          :class="executionStatus === 'success' ? 'text-green-500' : 'text-red-500'" class="text-4xl mb-4">
          <i :class="executionStatus === 'success' ? 'pi pi-check-circle' : 'pi pi-times-circle'"></i>
        </div>

        <!-- 状态标题 -->
        <h3 class="text-lg leading-6 font-medium mb-2" :class="{
          'text-blue-600': executionState === 'executing',
          'text-green-600': executionState === 'completed' && executionStatus === 'success',
          'text-red-600': executionState === 'completed' && executionStatus === 'error'
        }">
          <span v-if="executionState === 'executing'">命令执行中...</span>
          <span v-else-if="executionStatus === 'success'">命令执行成功</span>
          <span v-else>命令执行失败</span>
        </h3>

        <!-- 执行结果详情 -->
        <div class="mt-4 text-left">
          <!-- 总体执行状态 -->
          <div class="mb-4 p-3 rounded-md" :class="executionStatus === 'success' ? 'bg-green-50' : 'bg-red-50'">
            <h4 class="font-medium">执行状态:</h4>
            <p class="ml-2">{{ executionStatus === 'success' ? '成功' : '失败' }}</p>
          </div>

          <!-- 命令输出 (非技术性用户友好显示) -->
          <div
            v-if="commandOutputInternal !== undefined && commandOutputInternal !== null && commandOutputInternal !== ''"
            class="mb-4">
            <div class="flex justify-between items-center">
              <h4 class="font-medium">执行结果:</h4>
              <span class="text-sm px-2 py-1 bg-gray-100 rounded">{{ commandOutputInternal.length }} 字符</span>
            </div>
            <pre
              class="text-sm whitespace-pre-wrap p-3 bg-gray-100 rounded-md max-h-60 overflow-y-auto">{{ commandOutputInternal }}</pre>
          </div>

          <!-- 空输出提示 -->
          <div v-else class="mb-4">
            <h4 class="font-medium">执行结果:</h4>
            <div class="ml-2 p-3 bg-gray-50 rounded-md italic text-gray-600">命令执行完成，无输出内容</div>
          </div>

          <!-- 错误信息 (仅在有错误时显示) -->
          <div v-if="executionError && executionError !== ''" class="mb-4">
            <div class="flex justify-between items-center">
              <h4 class="font-medium text-red-600">错误信息:</h4>
              <span class="text-sm px-2 py-1 bg-red-100 text-red-800 rounded">{{ executionError.length }} 字符</span>
            </div>
            <pre
              class="text-sm whitespace-pre-wrap p-3 bg-red-100 text-red-800 rounded-md max-h-60 overflow-y-auto">{{ executionError }}</pre>
          </div>
          <!-- 不再显示空错误的提示 -->
        </div>

        <div class="items-center px-4 py-3">
          <button
            class="px-4 py-2 bg-blue-500 text-white text-base font-medium rounded-md w-full shadow-sm hover:bg-blue-600 focus:outline-none"
            @click="closeResultModal">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- 命令文本模态框 -->
  <div v-if="showCommandTextModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full"
    @click.self="showCommandTextModal = false">
    <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
      <div class="mt-3 text-center">
        <div class="text-gray-600 text-left">
          <h3 class="text-lg leading-6 font-medium text-gray-900">vars</h3>
          <pre class="text-sm whitespace-pre-wrap">{{ commandVariableValues }}</pre>
        </div>
        <h3 class="text-lg leading-6 font-medium text-gray-900">即将执行的命令</h3>
        <div class="mt-2 px-7 py-3">
          <p class="text-sm text-gray-500 break-all">{{ commandText }}</p>
        </div>
        <div class="items-center px-4 py-3">
          <button id="ok-btn"
            class="px-4 py-2 bg-blue-500 text-white text-base font-medium rounded-md w-full shadow-sm hover:bg-blue-600 focus:outline-none"
            @click="showCommandTextModal = false">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import { ExecuteCommand, GetCommandText } from '@/wailsjs/go/main/App';
import { useToastNotifications } from '@/composables/useToastNotifications';

const props = defineProps({
  selectedCommand: { type: Object as () => any, default: null },
  commandVariableValues: { type: Object as () => { [key: string]: any }, required: true },
  isProcessing: { type: Boolean, default: false },
  commandOutput: { type: String, default: '' },
  inputFilePath: { type: String, default: '' },
  outputFilePath: { type: String, default: '' },
});

const emit = defineEmits(['update:isProcessing', 'update:commandOutput']);

const { showToast } = useToastNotifications();

const isProcessingInternal = ref(props.isProcessing);
const commandOutputInternal = ref(props.commandOutput);

// 新增状态变量 - 更精确的状态管理
const showResultModal = ref(false);
const executionState = ref<'idle' | 'executing' | 'completed'>('idle'); // 'idle', 'executing', or 'completed'
const executionStatus = ref<'success' | 'error'>('success'); // 'success' or 'error'
const executionError = ref('');

watch(() => props.isProcessing, (newValue) => {
  isProcessingInternal.value = newValue;
});

watch(() => props.commandOutput, (newValue) => {
  commandOutputInternal.value = newValue;
});

watch(isProcessingInternal, (newValue) => {
  emit('update:isProcessing', newValue);
});

watch(commandOutputInternal, (newValue) => {
  emit('update:commandOutput', newValue);
});

const runCommand = async () => {
  if (!props.selectedCommand) {
    showToast('警告', '请选择要执行的命令', 'warn');
    return;
  }

  // 重置状态
  isProcessingInternal.value = true;
  commandOutputInternal.value = '';
  executionError.value = '';
  executionState.value = 'executing';
  executionStatus.value = 'success';
  showResultModal.value = true;

  try {
    const result = await ExecuteCommand(props.selectedCommand.id, props.commandVariableValues);
    commandOutputInternal.value = result;
    executionStatus.value = 'success';
  } catch (error) {
    executionStatus.value = 'error';
    executionError.value = String(error);
    commandOutputInternal.value = '';
  } finally {
    isProcessingInternal.value = false;
    executionState.value = 'completed';
  }
};

const closeResultModal = () => {
  showResultModal.value = false;
  // 重置状态当模态框关闭时
  executionState.value = 'idle';
  executionError.value = '';
};

const commandText = ref('');
const showCommandTextModal = ref(false);

const showCommandInfo = async () => {
  if (!props.selectedCommand) {
    showToast('警告', '请选择要执行的命令', 'warn');
    return;
  }

  try {
    const result = await GetCommandText(props.selectedCommand.id, props.commandVariableValues);
    commandText.value = result;
    showCommandTextModal.value = true;
  } catch (error) {
    showToast('错误', `获取命令文本失败: ${error}`, 'error');
    console.error('获取命令文本失败:', error);
  }
};
</script>