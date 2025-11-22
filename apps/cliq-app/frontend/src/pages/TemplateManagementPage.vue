<template>
  <div class="p-4">
    <div v-if="favTemplates.length === 0" class="text-center text-gray-500">
      <p>暂无收藏模板</p>
    </div>
    <div v-else>
      <DataTable :value="favTemplates" responsiveLayout="scroll">
        <Column field="name" header="模板名称"></Column>
        <Column field="description" header="描述"></Column>
        <Column field="author" header="作者"></Column>
        <Column field="version" header="版本"></Column>
        <Column header="操作">
          <template #body="slotProps">
            <Button icon="pi pi-pencil" size="small" @click="editTemplate(slotProps.data)" rounded variant="outlined" />
            <Button icon="pi pi-trash" size="small" @click="confirmDeleteTemplate(slotProps.data)" rounded
              variant="outlined" />
          </template>
        </Column>
      </DataTable>
    </div>

    <Dialog v-model:visible="displayConfirmation" header="确认删除" :modal="true">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span>您确定要删除模板 <b>{{ templateToDelete ? templateToDelete.name : '' }}</b> 吗？</span>
      </div>
      <template #footer>
        <Button label="取消" icon="pi pi-times" class="p-button-text" @click="closeConfirmation" />
        <Button label="删除" icon="pi pi-check" class="p-button-danger" @click="deleteTemplate" />
      </template>
    </Dialog>

    <!-- Template Editor Modal -->
    <TemplateEditorModal :visible="showEditorModal" :initialYaml="templateToEditContent"
      @close="showEditorModal = false" @save="onTemplateEdited" />
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { models } from '@/wailsjs/go/models';
import { ListFavTemplates, DeleteFavTemplate, GetFavTemplate, UpdateFavTemplate, GenerateYAMLFromTemplate, ParseYAMLToTemplate } from '@/wailsjs/go/main/App';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import { useToastNotifications } from '@/composables/useToastNotifications';
import TemplateEditorModal from '@/components/TemplateEditorModal.vue';

const favTemplates = ref<models.TemplateFile[]>([]);
const displayConfirmation = ref(false);
const templateToDelete = ref<models.TemplateFile | null>(null);

const showEditorModal = ref(false);
const templateToEdit = ref<models.TemplateFile | null>(null);
const templateToEditContent = ref('');
const { showToast } = useToastNotifications();

const loadFavTemplates = async () => {
  try {
    const result = await ListFavTemplates();
    favTemplates.value = result || [];
  } catch (error) {
    console.error('Failed to list favorite templates:', error);
    showToast('错误', `加载收藏模板失败: ${error}`, 'error');
  }
};

const confirmDeleteTemplate = (template: models.TemplateFile) => {
  templateToDelete.value = template;
  displayConfirmation.value = true;
};

const deleteTemplate = async () => {
  if (templateToDelete.value) {
    try {
      await DeleteFavTemplate(templateToDelete.value.name);
      showToast('成功', `模板 ${templateToDelete.value.name} 已删除`, 'success');
      await loadFavTemplates(); // Reload the list
    } catch (error) {
      console.error('Failed to delete template:', error);
      showToast('错误', `删除模板失败: ${error}`, 'error');
    } finally {
      closeConfirmation();
    }
  }
};

const editTemplate = async (template: models.TemplateFile) => {
  try {
    console.log('Attempting to edit template:', template.name);

    // Get the full template content
    const templateContent = await GetFavTemplate(template.name);
    console.log('Retrieved template content:', templateContent);

    if (!templateContent) {
      showToast('错误', `未找到模板内容: ${template.name}`, 'error');
      return;
    }

    // Convert the template object to YAML string
    const yamlContent = await GenerateYAMLFromTemplate(templateContent);
    console.log('Generated YAML content:', yamlContent);

    if (!yamlContent || yamlContent.trim() === '') {
      showToast('错误', `模板内容为空: ${template.name}`, 'error');
      return;
    }

    // Set the template to edit
    templateToEdit.value = templateContent;
    templateToEditContent.value = yamlContent;

    // Show the editor modal
    showEditorModal.value = true;
  } catch (error) {
    console.error('Failed to load template for editing:', error);
    showToast('错误', `加载模板进行编辑失败: ${error}`, 'error');
  }
};



const onTemplateEdited = async (updatedYaml: string) => {
  if (!templateToEdit.value) return;

  try {
    // Parse the updated YAML back to a template object
    const updatedTemplate = await ParseYAMLToTemplate(updatedYaml);

    // Send the updated template to the backend with old and new names
    await UpdateFavTemplate(templateToEdit.value.name, updatedTemplate.name, updatedTemplate);

    // Close the editor
    showEditorModal.value = false;

    // Reload the templates list
    await loadFavTemplates();

    showToast('成功', `模板 ${updatedTemplate.name} 已更新`, 'success');
  } catch (error) {
    console.error('Failed to update template:', error);
    showToast('错误', `更新模板失败: ${error}`, 'error');
  }
};

const closeConfirmation = () => {
  displayConfirmation.value = false;
  templateToDelete.value = null;
};

onMounted(() => {
  loadFavTemplates();
});
</script>

<style scoped>
.confirmation-content {
  display: flex;
  align-items: center;
}
</style>