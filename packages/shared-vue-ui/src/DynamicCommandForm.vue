<template>
  <div v-if="selectedCommand && commandVariables.length > 0" class="mb-6 p-4 bg-blue-50 rounded-md">
    <h3 class="font-medium mb-4">命令参数</h3>
    <div v-for="variable in commandVariables" :key="variable.name" class="mb-4">
      <label :for="variable.name" class="block text-sm font-medium text-gray-700 mb-2">
        {{ variable.label }}
        <span v-if="variable.required" class="text-red-500">*</span>
      </label>
      <!-- 文本输入 -->
      <InputText v-if="variable.type === 'string'" :id="variable.name"
        v-model="commandVariableValuesInternal[variable.name]" class="w-full" :placeholder="variable.description" />
      <!-- 数字输入 -->
      <InputNumber v-else-if="variable.type === 'number'" :id="variable.name"
        v-model="commandVariableValuesInternal[variable.name]" class="w-full" :placeholder="variable.description" />
      <!-- 布尔值 (Checkbox) -->
      <Checkbox v-else-if="variable.type === 'boolean'" :id="variable.name"
        v-model="commandVariableValuesInternal[variable.name]" :binary="true" />
      <!-- 文件输入 -->
      <div v-else-if="variable.type === 'file_input' || variable.type === 'file_output'"
        class="flex items-center space-x-2">
        <InputText :id="variable.name" v-model="commandVariableValuesInternal[variable.name]" readonly
          :placeholder="variable.description" class="w-full" />
        <Button type="button" @click="openFileSelection(variable.name, variable.type)" size="small" class="whitespace-nowrap">
          选择文件
        </Button>
      </div>
      <!-- 下拉选择 -->
      <Select v-else-if="variable.type === 'select'" :id="variable.name"
        v-model="commandVariableValuesInternal[variable.name]" :options="variable.options?.options || []"
        class="w-full" :placeholder="variable.description" />
      <small v-if="variable.description" class="mt-1 text-sm text-gray-500">{{ variable.description }}</small>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue';
import { OpenFileDialog, SaveFileDialog, OpenFileDialogWithFilters } from '@/wailsjs/go/main/App';
import { models } from '@/wailsjs/go/models';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import Checkbox from 'primevue/checkbox';
import { useToastNotifications } from '@/composables/useToastNotifications';

// Helper function to extract directory path and filename from a full path
function parsePath(inputPath: string): { directory: string; filename: string } {
  // Handle both Unix and Windows style paths
  const lastSlashIndex = Math.max(inputPath.lastIndexOf('/'), inputPath.lastIndexOf('\\'));
  
  if (lastSlashIndex === -1) {
    // No directory separator found, entire string is filename
    return { directory: '', filename: inputPath };
  } else {
    const directory = inputPath.substring(0, lastSlashIndex + 1);
    const filename = inputPath.substring(lastSlashIndex + 1);
    return { directory, filename };
  }
}

// Helper function to generate a new filename with "_new" suffix
function generateNewFilename(filename: string): string {
  const fileExtensionIndex = filename.lastIndexOf('.');

  if (fileExtensionIndex > 0) {
    // Split filename and extension (e.g., "document.pdf" -> "document" + ".pdf")
    const namePart = filename.substring(0, fileExtensionIndex);
    const extensionPart = filename.substring(fileExtensionIndex);
    return `${namePart}_new${extensionPart}`;
  } else {
    // No extension found, just add _new suffix (e.g., "README" -> "README_new")
    return `${filename}_new`;
  }
}

const props = defineProps({
  selectedCommand: { type: Object as () => any, default: null },
  commandVariableValues: { type: Object as () => { [key: string]: any }, required: true },
  inputFilePath: { type: String, default: '' },
  outputFilePath: { type: String, default: '' },
});

const emit = defineEmits(['update:commandVariableValues', 'update:inputFilePath', 'update:outputFilePath']);

const { showToast } = useToastNotifications();

const commandVariableValuesInternal = ref(props.commandVariableValues);
const inputFilePathInternal = ref(props.inputFilePath);
const outputFilePathInternal = ref(props.outputFilePath);

watch(() => props.commandVariableValues, (newValue) => {
  commandVariableValuesInternal.value = newValue;
}, { deep: true });


watch(commandVariableValuesInternal, (newValue) => {
  emit('update:commandVariableValues', newValue);
}, { deep: true });

// Watch for input file path changes to auto-populate output file path if not set
watch(inputFilePathInternal, (newInputPath) => {
  if (newInputPath) {
    // Find all output file variables in the command
    const outputFileVariables = commandVariables.value.filter(variable => variable.type === 'file_output');
    
    outputFileVariables.forEach(variable => {
      const currentOutputPath = commandVariableValuesInternal.value[variable.name];
      
      // If output path is not set, auto-populate it
      if (!currentOutputPath || currentOutputPath === '') {
        // Use helper functions to parse the path and generate new filename
        const { directory, filename } = parsePath(newInputPath);
        const newFileName = generateNewFilename(filename);
        
        // Create new output path by combining the directory path with new filename
        const outputFilePath = directory + newFileName;
        
        // Set the output file path
        commandVariableValuesInternal.value[variable.name] = outputFilePath;
        
        // If this is the main output file path prop, update that too
        if (variable.name === 'output_file' || variable.name === 'output') {
          outputFilePathInternal.value = outputFilePath;
          emit('update:outputFilePath', outputFilePath);
        }
      }
    });
  }
});


const commandVariables = computed(() => {
  if (props.selectedCommand && props.selectedCommand.variables) {
    // If variables is an array of VariableDefinition objects (with flattened structure)
    if (Array.isArray(props.selectedCommand.variables)) {
      return props.selectedCommand.variables.map((varDef: models.VariableDefinition) => ({
        name: varDef.name,
        type: varDef.type,
        arg_name: varDef.arg_name,
        label: varDef.label,
        description: varDef.description,
        required: varDef.required,
        options: varDef.options,
      }));
    } else {
      // Fallback for old map format (for backward compatibility if needed)
      return Object.entries(props.selectedCommand.variables).map(([name, variable]) => ({
        name,
        ...(variable as models.Variable),
      }));
    }
  }
  return [];
});

const openFileSelection = async (variableName: string, variableType: string) => {
  let filePath = '';
  try {
    let variableDef: models.VariableDefinition | undefined;
    
    // Find the variable in the array format
    if (Array.isArray(props.selectedCommand.variables)) {
      variableDef = (props.selectedCommand.variables as models.VariableDefinition[])
        .find(v => v.name === variableName);
    } else {
      // Fallback for old map format
      const varMap = props.selectedCommand.variables as Record<string, models.Variable>;
      if (varMap && varMap[variableName]) {
        const variable = varMap[variableName];
        // Create a compatible structure for old format
        variableDef = {
          name: variableName,
          type: variable.type,
          arg_name: variable.arg_name,
          label: variable.label,
          description: variable.description,
          required: variable.required,
          options: variable.options,
        } as models.VariableDefinition;
      }
    }

    if (variableDef) {
      if (variableType === 'file_input') {
        if (variableDef.options && 
            variableDef.options.file_types && 
            Array.isArray(variableDef.options.file_types) && 
            variableDef.options.file_types.length > 0) {
          // Use specific file type filters
          filePath = await OpenFileDialogWithFilters(formatFileFilters(variableDef.options.file_types));
        } else {
          // Fallback to all files when no file_types are specified or are empty
          filePath = await OpenFileDialog();
        }
        if (filePath) {
          inputFilePathInternal.value = filePath;
        }
      } else if (variableType === 'file_output') {
        filePath = await SaveFileDialog();
        if (filePath) {
          outputFilePathInternal.value = filePath;
        }
      }
    } else {
      // Fallback for old behavior
      if (variableType === 'file_input') {
        filePath = await OpenFileDialog();
        if (filePath) {
          inputFilePathInternal.value = filePath;
        }
      } else if (variableType === 'file_output') {
        filePath = await SaveFileDialog();
        if (filePath) {
          outputFilePathInternal.value = filePath;
        }
      }
    }

    if (filePath) {
      commandVariableValuesInternal.value[variableName] = filePath;
    }
  } catch (error) {
    showToast('错误', `选择文件失败: ${error}`, 'error');
    console.error('选择文件失败:', error);
  }
};

/**
 * Function to format file type filters for the backend.
 * 
 * Expected format for fileTypes: array of file extensions, e.g. ['.pdf', 'docx', 'xlsx']
 * Each entry will be normalized to start with a dot and have no leading/trailing whitespace.
 * Wildcards are automatically handled as '*.<ext>'.
 */
const formatFileFilters = (fileTypes: string[]) => {
  if (!fileTypes || !Array.isArray(fileTypes) || fileTypes.length === 0) {
    return [];
  }

  // Normalize file types: ensure each starts with a dot, remove whitespace, and remove any leading wildcards
  const normalizedTypes = fileTypes.map(type => {
    let ext = type.trim().replace(/^\*+/, ''); // Remove leading wildcards
    if (!ext.startsWith('.')) {
      ext = '.' + ext.replace(/^\./, '');
    }
    return ext;
  });

  // Join the normalized file types with semicolons for the pattern
  const pattern = normalizedTypes.map(ext => `*${ext}`).join(';');

  // Create a display name from the normalized file types
  const displayName = `支持的文件 (${normalizedTypes.join(', ')})`;

  return [{
    DisplayName: displayName,
    Pattern: pattern
  }];
};
</script>