<template>
  <div v-if="template" class="max-w-4xl mx-auto">
    <Card>
      <template #title>
        <div class="flex justify-between items-start">
          <h1 class="text-3xl font-bold">{{ template.title }}</h1>
          <Button label="Download" icon="pi pi-download" @click="downloadTemplate" />
        </div>
      </template>
      <template #subtitle>
        <div class="flex gap-4 items-center">
          <span>By {{ template.author?.username }}</span>
          <span>{{ formatDate(template.CreatedAt) }}</span>
          <span>{{ template.downloads }} downloads</span>
        </div>
      </template>
      <template #content>
        <div class="mt-6">
          <h3 class="text-xl font-semibold mb-2">Description</h3>
          <p class="text-gray-700 whitespace-pre-wrap">{{ template.description }}</p>
        </div>

        <div class="mt-8">
          <h3 class="text-xl font-semibold mb-2">Preview</h3>
          <div class="bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto font-mono text-sm">
            <pre>{{ template.content }}</pre>
          </div>
        </div>
      </template>
    </Card>
  </div>
  <div v-else class="text-center py-12">
    Loading...
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import { useTemplateStore } from '../stores/template';
import Card from 'primevue/card';
import Button from 'primevue/button';

const route = useRoute();
const templateStore = useTemplateStore();

const template = computed(() => templateStore.currentTemplate);

onMounted(async () => {
  const id = route.params.id as string;
  await templateStore.fetchTemplate(id);
});

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString();
};

const downloadTemplate = () => {
  if (!template.value) return;
  const blob = new Blob([template.value.content], { type: 'text/yaml' });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${template.value.title.replace(/\s+/g, '_').toLowerCase()}.cliqfile.yaml`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  window.URL.revokeObjectURL(url);
};
</script>
